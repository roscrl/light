package app

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/roscrl/light/config"
	"github.com/roscrl/light/core/helpers/ulid"
	"github.com/roscrl/light/db"
)

func BenchmarkHandleTodosCreate(b *testing.B) {
	b.ReportAllocs()

	tempDBPath := db.NewTempDBFileBenchmarkWithCleanup(b)

	cfg := config.NewTestConfig()
	cfg.SqliteDBPath = tempDBPath

	app := NewAppBenchmarkWithCleanup(b, cfg)
	db.RunMigrations(app.DB)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()

			body := fmt.Sprintf("task=%s", ulid.NewString())

			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
			r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			app.handleTodosCreate()(w, r)

			if w.Result().StatusCode != http.StatusOK {
				b.Error("expected status code 200")
			}
		}
	})
}
