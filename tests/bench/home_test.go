package bench

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/roscrl/light/app"
	"github.com/roscrl/light/config"
	"github.com/roscrl/light/core/helpers/ulid"
	"github.com/roscrl/light/db"
)

func BenchmarkHome(b *testing.B) {
	b.ReportAllocs()

	cfg := config.NewTestConfig()
	cfg.SqliteDBPath = fmt.Sprintf("file:%s?mode=memory&cache=shared", ulid.NewString())

	app := app.NewApp(cfg)

	b.Cleanup(func() {
		app.DB.Close()
	})

	db.RunMigrations(app.DB, db.PathMigrations)

	_, _ = app.DB.Exec("INSERT INTO todos (id, task, status) VALUES (?, ?, ?)", ulid.NewString(), "important todo!", "pending")
	_, _ = app.DB.Exec("INSERT INTO todos (id, task, status) VALUES (?, ?, ?)", ulid.NewString(), "also important!", "done")

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			app.ServeHTTP(w, r)

			if w.Result().StatusCode != http.StatusOK {
				b.Error("expected status code 200")
			}
		}
	})
}
