package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akasrt/filensy/internal/config/env"
	"github.com/akasrt/filensy/internal/database"
	"github.com/akasrt/filensy/internal/file"
	"github.com/akasrt/filensy/internal/filestore"
	"github.com/akasrt/filensy/internal/middlewarex"
	"github.com/akasrt/filensy/internal/router"
	"github.com/akasrt/filensy/internal/util/errorx"
	"github.com/akasrt/filensy/internal/util/loggerx"
	"github.com/akasrt/filensy/internal/util/validate"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func Start() {
	initDependencies()
	e := echo.New()
	e.Use(middleware.RequestID())
	e.Use(middlewarex.LoggerMiddleware())
	e.Validator = validate.NewValidator()
	e.HTTPErrorHandler = errorx.ErrorHandler()
	router.AddRoutes(e)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	address := env.GetEnv(env.ServerAddress)
	sc := echo.StartConfig{
		Address:         address,
		GracefulTimeout: 10 * time.Second,
		HideBanner:      true,
	}

	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}

	cleanUp()
}

func initDependencies() {
	env.LoadEnv()
	loggerx.Init()
	database.InitDB()
	filestore.InitFileStore(env.GetEnv(env.FileRoot))
	file.RunCleaner()
}

func cleanUp() {
	database.Close()
}
