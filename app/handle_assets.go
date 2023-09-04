package app

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/roscrl/light/config"
)

const (
	PathAssets         = "core/views/assets/dist"
	PathEmbeddedAssets = "assets/dist"
)

func (app *App) handleAssets() http.HandlerFunc {
	if app.Cfg.Env == config.LOCAL {
		assetsFileServer := http.FileServer(http.Dir("./" + PathAssets + "/"))
		handler := http.StripPrefix(RouteAssetsBase+"/", assetsFileServer)

		return handler.ServeHTTP
	}

	subFS, err := fs.Sub(app.Cfg.FrontendDistFS, PathEmbeddedAssets)
	if err != nil {
		log.Fatal(err)
	}

	assetFileServer := http.FileServer(http.FS(subFS))
	handler := http.StripPrefix(RouteAssetsBase+"/", assetFileServer)

	return handler.ServeHTTP
}
