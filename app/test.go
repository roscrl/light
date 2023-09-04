package app

import (
	"testing"

	"github.com/matryer/is"
	"github.com/roscrl/light/config"
)

func NewUnstartedTestApp(t *testing.T, cfg *config.App) (*is.I, *App) {
	t.Helper()

	is, app := is.New(t), NewApp(cfg)

	return is, app
}

func NewStartedTestAppWithCleanup(t *testing.T, cfg *config.App) (*is.I, *App) {
	t.Helper()

	is, app := is.New(t), NewApp(cfg)
	is.NoErr(app.Start())

	t.Cleanup(func() {
		is.NoErr(app.Stop())
	})

	return is, app
}
