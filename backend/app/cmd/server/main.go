package main

import (
	"backend/internal/app"
	"backend/internal/config"
	"backend/pkg/logging"
	"log"
)

func main() {
	log.Print("config init")
	cfg := config.GetConfig()

	log.Print("logger init")
	logger := logging.GetLogger(cfg.AppConfig.LogLevel)

	a, err := app.NewApp(cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Println("Running Application")
	a.Run()
}
