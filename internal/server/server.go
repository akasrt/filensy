package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akasrt/filensy/internal/config/env"
	"github.com/akasrt/filensy/internal/database"
	"github.com/labstack/echo/v5"
)

func Start() {
	initDependencies()
	e := echo.New()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	address := ":" + env.GetEnv(env.ServerPort)
	sc := echo.StartConfig{
		Address:         address,
		GracefulTimeout: 10 * time.Second,
		HideBanner:      true,
		OnShutdownError: cleanUp,
	}

	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}

func initDependencies() {
	database.InitDB()
}

func cleanUp(err error) {
	database.Close()
}
