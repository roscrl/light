package core

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/roscrl/light/config"
)

//go:embed views/assets/dist
var assetsFS embed.FS

const (
	PathAssets         = "core/views/assets/dist"
	PathEmbeddedAssets = "views/assets/dist"
)

func (s *Server) handleAssets() http.HandlerFunc {
	if s.Cfg.Env == config.DEV {
		assetsFileServer := http.FileServer(http.Dir("./" + PathAssets + "/"))
		handler := http.StripPrefix(RouteAssetBase+"/", assetsFileServer)

		return handler.ServeHTTP
	}

	subFS, err := fs.Sub(assetsFS, PathEmbeddedAssets)
	if err != nil {
		log.Fatal(err)
	}

	assetFileServer := http.FileServer(http.FS(subFS))
	handler := http.StripPrefix(RouteAssetBase+"/", assetFileServer)

	return handler.ServeHTTP
}
