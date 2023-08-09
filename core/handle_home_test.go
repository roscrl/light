package core

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
	"github.com/roscrl/light/config"
	_ "github.com/roscrl/light/core/support/testhelp"
)

func TestHandleHome(t *testing.T) {
	is, server := is.New(t), NewServer(config.TestConfig())

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	server.ServeHTTP(w, r) // integration test like (middlewares included)
	is.Equal(w.Result().StatusCode, http.StatusOK)

	server.handleHome()(w, r) // unit test like (no middlewares)
	is.Equal(w.Result().StatusCode, http.StatusOK)
}

//func TestHandleTopPlaylistsAfterCursor(t *testing.T) {
//	mc := config.TestConfig()
//	mc.SqliteDBPath = ":memory:"
//
//	is, server := is.New(t), NewServer(mc)
//
//	db.RunMigrations(server.DB, db.PathMigrations)
//
//	w := httptest.NewRecorder()
//
//	req := httptest.NewRequest(http.MethodGet, "/playlists/top?after=6-10", nil)
//	req.Header.Set("Accept", views.TurboStreamMIME)
//
//	server.ServeHTTP(w, req)
//
//	server.handlePlaylistsPaginationTop()(w, req)
//	is.Equal(w.Result().StatusCode, http.StatusOK)
//}
