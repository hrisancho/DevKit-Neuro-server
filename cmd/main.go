package main

import (
	"DevKit-Neuro-server/internal/config"
	main_logger "DevKit-Neuro-server/internal/logger"
	"DevKit-Neuro-server/internal/server"
	main_validator "DevKit-Neuro-server/internal/validator"
	"context"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, ctxCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer ctxCancel()

	mainValidator, err := main_validator.NewValidator()

	mainConfig, err := config.LoadConfig(mainValidator)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := main_logger.NewLogger(config.DebugLevel)
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("config loaded", zap.Any("config", mainConfig))

	mainServer, err := server.NewServer(mainConfig, logger, mainValidator)
	if err != nil {
		logger.Fatal("error while starting server", zap.Error(err))
	}

	defer mainServer.App.Close()

	mainServer.Run(ctx)

}
