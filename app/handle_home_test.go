package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/roscrl/light/config"
	_ "github.com/roscrl/light/core/utils/testutil"
)

func TestHandleHome(t *testing.T) {
	t.Parallel()

	is, app := NewUnstartedTestApp(t, config.NewTestConfig())

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, RouteHome, nil)

	app.ServeHTTP(w, r) // integration test like (middlewares included)
	is.Equal(w.Result().StatusCode, http.StatusOK)

	app.handleHome()(w, r) // unit test like (no middlewares)
	is.Equal(w.Result().StatusCode, http.StatusOK)
}
