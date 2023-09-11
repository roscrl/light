package jobs

import (
	"context"
	"log"
	"log/slog"
	"testing"
	"time"

	"github.com/matryer/is"

	"github.com/roscrl/light/db"
)

func TestProcessorOneJob(t *testing.T) {
	t.Parallel()

	is := is.New(t)
	_, qry := db.NewTempMigratedDBAndQueriesTestingWithCleanup(t)

	testJob := JobName("test")

	processor := &Processor{
		Qry:         qry,
		Interval:    time.Millisecond * 100,
		Log:         slog.Default(),
		JobFinished: make(chan string),
		JobNameToJobFuncRegistry: JobNameToJobFuncRegistry{
			testJob: func(args map[string]any) error {
				log.Println("test job!")

				return nil
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go processor.StartJobLoop(ctx)

	enqueuedJobID, err := Enqueue(ctx, testJob, nil, qry)
	is.NoErr(err)

	finishedJobID := <-processor.JobFinished
	is.Equal(enqueuedJobID, finishedJobID)
}

func TestProcessorMultipleJobs(t *testing.T) {
	t.Parallel()

	is := is.New(t)
	_, qry := db.NewTempMigratedDBAndQueriesTestingWithCleanup(t)

	testJob := JobName("test")
	testJob2 := JobName("test2")

	processor := &Processor{
		Qry:         qry,
		Interval:    time.Millisecond * 100,
		Log:         slog.Default(),
		JobFinished: make(chan string),
		JobNameToJobFuncRegistry: JobNameToJobFuncRegistry{
			testJob: func(args map[string]any) error {
				is.Equal(args["hello"], "world!")

				log.Println("test job!", args["hello"])

				return nil
			},
			testJob2: func(args map[string]any) error {
				log.Println("test job 2!")

				return nil
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go processor.StartJobLoop(ctx)

	enqueuedJobID1, err := Enqueue(ctx, testJob, map[string]any{
		"hello": "world!",
	}, qry)
	is.NoErr(err)

	enqueuedJobID2, err := Enqueue(ctx, testJob2, nil, qry)
	is.NoErr(err)

	//nolint:gosimple
	for i := 0; i < 2; i++ {
		select {
		case finishedJobID := <-processor.JobFinished:
			switch {
			case finishedJobID == enqueuedJobID1:
				t.Log("finished job 1")
			case finishedJobID == enqueuedJobID2:
				t.Log("finished job 2")
			default:
				t.Fatal("finished job ID did not match any enqueued job ID")
			}
		}
	}
}

func TestProcessorPanic(t *testing.T) {
	t.Parallel()

	is := is.New(t)
	_, qry := db.NewTempMigratedDBAndQueriesTestingWithCleanup(t)

	testJob := JobName("test")

	processor := &Processor{
		Qry:         qry,
		Interval:    time.Millisecond * 100,
		Log:         slog.Default(),
		JobFinished: make(chan string),
		JobNameToJobFuncRegistry: JobNameToJobFuncRegistry{
			testJob: func(args map[string]any) error {
				panic("job panic!")
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go processor.StartJobLoop(ctx)

	enqueuedJobID, err := Enqueue(ctx, testJob, nil, qry)
	is.NoErr(err)

	finishedJobID := <-processor.JobFinished
	is.Equal(enqueuedJobID, finishedJobID)
}
