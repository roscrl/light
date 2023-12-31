package bench

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/roscrl/light/app"
	"github.com/roscrl/light/config"
	"github.com/roscrl/light/core/helpers/ulid"
	"github.com/roscrl/light/db"

	_ "github.com/roscrl/light/core/utils/testutil"
)

func BenchmarkHome(b *testing.B) {
	b.ReportAllocs()

	tempDBPath := db.NewTempDBFileBenchmarkWithCleanup(b)

	cfg := config.NewTestConfig()
	cfg.SqliteDBPath = tempDBPath

	app := app.NewAppBenchmarkWithCleanup(b, cfg)
	db.RunMigrations(app.DB)

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
