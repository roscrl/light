package core

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/roscrl/light/config"
	"github.com/roscrl/light/core/app"
)

func Bootstrap() {
	var configPath string

	flag.StringVar(&configPath, "config", "USE_ENVIRONMENT", "file path to server config file else use environment variables")
	flag.Parse()

	var cfg *config.App
	if configPath == "USE_ENVIRONMENT" {
		cfg = config.NewFromEnv()
	} else {
		cfg = config.NewFromCustomConfig(configPath)
	}

	app := app.NewApp(cfg)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	if err := app.Start(); err != nil {
		log.Fatalf("failed to start app: %v", err)
	}

	<-stop

	if err := app.Stop(); err != nil {
		log.Fatalf("failed to stop app: %v", err)
	}
}
