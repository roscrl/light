package core

import (
	"embed"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/roscrl/light/config"
)

//go:embed views/assets/dist
var frontendAssets embed.FS

func Bootstrap() {
	var configPath string

	flag.StringVar(&configPath, "config", "USE_ENVIRONMENT", "file path to server config file otherwise use environment variables")
	flag.Parse()

	var cfg *config.Server
	if configPath == "USE_ENVIRONMENT" {
		cfg = config.FromEnv()
	} else {
		cfg = config.FromCustomConfig(configPath)
	}

	cfg.FrontendAssetsFS = frontendAssets
	cfg.MustValidate()

	srv := NewServer(cfg)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	srv.Start()

	<-stop

	srv.Stop()
}
