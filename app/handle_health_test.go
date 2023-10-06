package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/roscrl/light/config"

	_ "github.com/roscrl/light/core/utils/testutil"
)

func TestHandleHealth(t *testing.T) {
	t.Parallel()

	is, app := NewAppUnstarted(t, config.NewTestConfig())

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, RouteHealth, nil)

	app.ServeHTTP(w, r)
	is.Equal(w.Result().StatusCode, http.StatusOK)
}
