package jobs

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/roscrl/light/core/helpers/rlog/key"
	"github.com/roscrl/light/core/helpers/rlog/keygroup"
	"github.com/roscrl/light/core/jobs/scope"
	"github.com/roscrl/light/db/sqlc"
)

// Processor is a simple job processor that fetches past due jobs
// from the database and runs them. It is expected that there is
// only one Processor singleton in the app and across any future
// horizontal scaling of the app. This is because the Processor
// does not lock the jobs table when fetching due jobs, so it is
// possible that two processors fetch the same job and run it twice.
// For more complex job processing needs, consider using a dedicated
// job processing library. See Enqueue and Schedule for adding jobs.
type Processor struct {
	Qry *sqlc.Queries

	Interval time.Duration
	Log      *slog.Logger

	JobNameToJobFuncRegistry JobNameToJobFuncRegistry

	JobsInFlight sync.WaitGroup

	// JobFinished is a channel that is sent the ID of a job
	// when it is finished processing. This is useful for testing
	// purposes. Finished is defined as the job being set to
	// success or failed in the database or failed to set to
	// running in the database.
	JobFinished chan string
}

// StartJobLoop initiates the job processing loop that periodically
// checks and processes due jobs. A job is considered due if its run_at
// column is in the past and its status is pending. The loop runs every
// Interval duration. So there is a possibility that a job is not run
// at the exact time it is due, but it will be run on the next loop.
func (p *Processor) StartJobLoop(ctx context.Context, jobScope *scope.Job) {
	ticker := time.NewTicker(p.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := p.processDueJobs(ctx, jobScope)
			if err != nil {
				p.Log.Error("processing due jobs", key.Err, err)
			}
		case <-ctx.Done():
			p.Log.InfoContext(ctx, "context done, waiting for any remaining jobs to finish", "ctx", ctx.Err())

			p.JobsInFlight.Wait()

			p.Log.InfoContext(ctx, "all jobs finished, exiting job loop")

			return
		}
	}
}

func (p *Processor) processDueJobs(ctx context.Context, jobScope *scope.Job) error {
	jobs, err := p.Qry.GetOverduePendingJobsFromTime(ctx, time.Now().Unix())
	if err != nil {
		return fmt.Errorf("getting db pending jobs: %w", err)
	}

	for _, job := range jobs {
		job := job

		log := p.Log.WithGroup(keygroup.Job)
		log = log.With(key.ID, job.ID, key.Name, job.Name, key.RunAt, job.RunAt)

		jobFunc := p.JobNameToJobFuncRegistry[JobName(job.Name)]
		if jobFunc == nil {
			log.ErrorContext(ctx, "attempted to run due job but no matching job function found")

			failedJobParams := sqlc.SetFailedJobParams{
				FailedMessage: sql.NullString{
					String: err.Error(),
					Valid:  true,
				},
				ID: job.ID,
			}

			if err = p.Qry.SetFailedJob(ctx, failedJobParams); err != nil {
				log.ErrorContext(ctx, "setting db job status to failed with failure message", key.Err, err)
			}

			continue
		}

		var args map[string]any
		if job.Arguments != "" {
			if err := json.Unmarshal([]byte(job.Arguments), &args); err != nil {
				log.ErrorContext(ctx, "unmarshalling job arguments", key.Err, err)

				failedJobParams := sqlc.SetFailedJobParams{
					FailedMessage: sql.NullString{
						String: err.Error(),
						Valid:  true,
					},
					ID: job.ID,
				}

				if err = p.Qry.SetFailedJob(ctx, failedJobParams); err != nil {
					log.Error("setting db job status to failed with failure message", key.Err, err)
				}

				continue
			}
		}

		p.JobsInFlight.Add(1)

		go func() {
			defer func() {
				p.JobsInFlight.Done()

				select {
				case p.JobFinished <- job.ID:
				default:
				} // Don't worry if JobID can't be sent.

				if recovery := recover(); recovery != nil {
					var err error
					switch panicType := recovery.(type) {
					case string:
						err = fmt.Errorf(panicType)
					case error:
						err = panicType
					default:
						err = fmt.Errorf("unknown panic: %v", panicType)
					}

					log.ErrorContext(ctx, "panic during job processing", key.Err, err)

					failedJobParams := sqlc.SetFailedJobParams{
						FailedMessage: sql.NullString{
							String: err.Error(),
							Valid:  true,
						},
						ID: job.ID,
					}

					if err := p.Qry.SetFailedJob(ctx, failedJobParams); err != nil {
						log.ErrorContext(ctx, "setting db job status to failed with failure message", key.Err, err)

						return
					}
				}
			}()

			if err := p.Qry.SetJobStatusToRunning(ctx, job.ID); err != nil {
				log.ErrorContext(ctx, "setting db job status to running", key.Err, err)

				return
			}

			err := jobFunc(ctx, jobScope, args)
			if err != nil {
				log.ErrorContext(ctx, "job failed", key.Err, err)

				failedJobParams := sqlc.SetFailedJobParams{
					FailedMessage: sql.NullString{
						String: err.Error(),
						Valid:  true,
					},
					ID: job.ID,
				}

				if err := p.Qry.SetFailedJob(ctx, failedJobParams); err != nil {
					log.ErrorContext(ctx, "setting db job status to failed with failure message", key.Err, err)
				}

				return
			}

			if err := p.Qry.SetSuccessfulJob(ctx, job.ID); err != nil {
				log.ErrorContext(ctx, "setting db job status to success", key.Err, err)

				return
			}
		}()
	}

	return nil
}
