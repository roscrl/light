package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/roscrl/light/config"
	"github.com/roscrl/light/core/helpers/rlog"
	"github.com/roscrl/light/core/helpers/rlog/key"
	"github.com/roscrl/light/core/views"
	"github.com/roscrl/light/db"
	"github.com/roscrl/light/db/sqlc"
)

type App struct {
	Cfg *config.App
	Log *slog.Logger

	DB  *sql.DB
	Qry *sqlc.Queries

	Views *views.Views

	Client *http.Client

	Router   http.Handler
	Listener net.Listener
	Port     string

	HTTPServer *http.Server
}

//nolint:gomnd
func NewApp(cfg *config.App) *App {
	cfg.FrontendDistFS = views.FrontendDistFS
	cfg.MustValidate()

	app := &App{}

	app.Cfg = cfg
	app.Log = rlog.NewDefaultLogger()
	slog.SetDefault(app.Log)

	app.DB = db.New(cfg.SqliteDBPath)
	app.Qry = sqlc.New(app.DB)

	app.Views = views.New(app.Cfg.Env)

	app.Client = &http.Client{
		Timeout: 10 * time.Second,
	}

	app.services()

	app.Router = app.routes()

	app.HTTPServer = &http.Server{
		Handler:      app.Router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 35 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return app
}

func (app *App) Start() error {
	app.Log.Info("starting server", key.Env, app.Cfg.Env, key.DB, app.Cfg.SqliteDBPath)

	listener, err := net.Listen("tcp", ":"+app.Cfg.Port)
	if err != nil {
		return fmt.Errorf("listening on port %v: %w", app.Cfg.Port, err)
	}

	app.Listener = listener
	app.Port = fmt.Sprintf("%v", listener.Addr().(*net.TCPAddr).Port)
	app.Cfg.Port = app.Port

	go func() {
		err := app.HTTPServer.Serve(app.Listener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			var opErr *net.OpError
			if errors.As(err, &opErr) && opErr.Op == "accept" {
				app.Log.Info("server shut down")
			} else {
				log.Fatal("failed to stop server: ", err)
			}
		}
	}()

	app.Log.Info("ready to handle requests on port " + app.Port)

	return nil
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.Router.ServeHTTP(w, r)
}

func (app *App) Stop() error {
	app.Log.Info("server shutting down...")

	if err := app.Listener.Close(); err != nil {
		return fmt.Errorf("closing listener: %w", err)
	}

	app.Views.StopLocalBrowserRefreshChannelIfLocal()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := app.HTTPServer.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("shutting down http server: %w", err)
	}

	err = app.DB.Close()
	if err != nil {
		return fmt.Errorf("closing database: %w", err)
	}

	return nil
}