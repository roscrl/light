package app

import (
	"net/http"
	"time"

	"github.com/roscrl/light/core/helpers/rlog"
	"github.com/roscrl/light/core/helpers/rlog/key"
	"github.com/roscrl/light/core/views"
	"github.com/roscrl/light/core/views/params"
	"github.com/roscrl/light/db/sqlc"
)

const jobFetchLimit = 1_000

func (app *App) handleAdmin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log, rctx := rlog.L(r)

		jobs, err := app.Qry.GetJobs(rctx, jobFetchLimit)
		if err != nil {
			log.ErrorContext(rctx, "failed to query for jobs", key.Err, err)
			app.Views.RenderDefaultErrorPage(w)

			return
		}

		app.Views.RenderPage(w, views.Admin, map[string]any{
			params.Jobs: jobs,
		})
	}
}

func (app *App) handleAdminJobSearch() http.HandlerFunc {
	const (
		fieldID           = "id"
		fieldRunOnOrAfter = "run_on_or_after"
	)

	return func(w http.ResponseWriter, r *http.Request) {
		log, rctx := rlog.L(r)

		var (
			jobs []sqlc.Job
			job  sqlc.Job
			err  error
		)

		id, runAtFrom := r.FormValue(fieldID), r.FormValue(fieldRunOnOrAfter)

		givenRunAt := runAtFrom != ""
		givenID := id != ""

		switch {
		case givenRunAt:
			var givenTimeParsed time.Time
			{
				givenTimeParsed, err = time.Parse("2006-01-02", runAtFrom)
				if err != nil {
					log.ErrorContext(rctx, "failed to get by run at time", key.Err, err)
					app.Views.RenderTurboStream(w, views.JobFormSearchStream, map[string]any{
						params.InputJob: id,
						params.Error:    "Something went wrong parsing the run at time",
					})

					return
				}
			}

			givenTimeAsUnix := givenTimeParsed.Unix()
			jobs, err = app.Qry.GetJobsAfterEqualRunAtDate(rctx, sqlc.GetJobsAfterEqualRunAtDateParams{
				From:  givenTimeAsUnix,
				Limit: jobFetchLimit,
			})
		case givenID:
			job, err = app.Qry.GetJobByID(rctx, id)
			jobs = append(jobs, job)
		default:
			jobs, err = app.Qry.GetJobs(rctx, jobFetchLimit)
		}

		if err != nil {
			log.ErrorContext(rctx, "failed to get by job id", key.Err, err)
			app.Views.RenderTurboStream(w, views.JobFormSearchStream, map[string]any{
				params.InputJob: id,
				params.Error:    "Oops, something went wrong searching by job id",
			})

			return
		}

		app.Views.RenderPage(w, views.Admin, map[string]any{
			params.InputJob: id,
			params.Jobs:     jobs,
		})
	}
}
