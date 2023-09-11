package app

import (
	"context"
	"time"

	"github.com/roscrl/light/core/jobs"
)

//nolint:revive,staticcheck,gocritic,wsl
func (app *App) services(ctx context.Context) {
	app.JobsProcessor = &jobs.Processor{
		Qry:                      app.Qry,
		Interval:                 time.Second * 5,
		Log:                      app.Log,
		JobNameToJobFuncRegistry: jobs.DefaultRegistry(),
	}

	go app.JobsProcessor.StartJobLoop(ctx)

	if app.Cfg.Mocking {
		// mockedServices
	} else {
		// realServices
	}

	// sharedServices
}
