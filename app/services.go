package app

import (
	"context"
	"time"

	"github.com/roscrl/light/core/jobs"
	"github.com/roscrl/light/core/jobs/scope"
)

//nolint:revive,staticcheck,gocritic,wsl
func (app *App) services(ctx context.Context) {
	app.JobsProcessor = &jobs.Processor{
		Qry:                      app.Qry,
		Interval:                 time.Second * 5,
		Log:                      app.Log,
		JobNameToJobFuncRegistry: jobs.DefaultRegistry(),
	}

	{
		jobScope := &scope.Job{
			Cfg:    app.Cfg,
			DB:     app.DB,
			Qry:    app.Qry,
			Client: app.Client,
		}

		go app.JobsProcessor.StartJobLoop(ctx, jobScope)
	}

	if app.Cfg.Mocking {
		// mockedServices
	} else {
		// realServices
	}

	// sharedServices
}
