package app

import (
	"context"
	"time"

	"github.com/roscrl/light/core/jobs"
)

//nolint:revive,staticcheck,gocritic,wsl
func (app *App) services(ctx context.Context) {
	app.Jobs = &jobs.Processor{
		Qry:         app.Qry,
		Interval:    time.Second * 5,
		Log:         app.Log,
		JobRegistry: jobs.DefaultRegistry(),
	}

	go app.Jobs.StartJobLoop(ctx)

	if app.Cfg.Mocking {
		// mockedServices
	} else {
		// realServices
	}

	// sharedServices
}
