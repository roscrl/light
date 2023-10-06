package app

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/roscrl/light/config"
	"github.com/roscrl/light/core/helpers/ulid"
	"github.com/roscrl/light/db"

	_ "github.com/roscrl/light/core/utils/testutil"
)

func TestHandleHome(t *testing.T) {
	t.Parallel()

	cfg := config.NewTestConfig()
	cfg.SqliteDBPath = fmt.Sprintf("file:%s?mode=memory", ulid.NewString())

	is, app := NewAppStartedWithCleanup(t, cfg)
	db.RunMigrations(app.DB)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, RouteHome, nil)

	app.ServeHTTP(w, r) // integration test like (middlewares included)
	is.Equal(w.Result().StatusCode, http.StatusOK)

	app.handleHome()(w, r) // unit test like (no middlewares)
	is.Equal(w.Result().StatusCode, http.StatusOK)
}
