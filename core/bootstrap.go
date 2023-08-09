package core

import (
	"flag"
	"github.com/roscrl/light/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

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

	srv := NewServer(cfg)
	slog.SetDefault(srv.Log)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	srv.Start()

	<-stop

	srv.Stop()
}
