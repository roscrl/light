package app

import (
	"context"
	"testing"

	"github.com/matryer/is"

	"github.com/roscrl/light/config"
)

func NewUnstartedTestApp(t *testing.T, cfg *config.App) (*is.I, *App) {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	is, app := is.New(t), NewApp(ctx, cfg)

	t.Cleanup(func() {
		cancel()
	})

	return is, app
}

func NewAppBenchmarkWithCleanup(b *testing.B, cfg *config.App) *App {
	ctx, cancel := context.WithCancel(context.Background())
	app := NewApp(ctx, cfg)

	b.Cleanup(func() {
		cancel()
		err := app.DB.Close()
		if err != nil {
			b.Fatal(err)
		}
	})

	return app
}

func NewStartedAppWithCleanup(t *testing.T, cfg *config.App) (*is.I, *App) {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())

	is, app := is.New(t), NewApp(ctx, cfg)
	is.NoErr(app.Start())

	t.Cleanup(func() {
		cancel()
		is.NoErr(app.Stop())
	})

	return is, app
}
