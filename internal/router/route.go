package router

import (
	"github.com/akasrt/filensy/internal/file"
	"github.com/akasrt/filensy/internal/middlewarex"
	"github.com/labstack/echo/v5"
)

func AddRoutes(e *echo.Echo) {
	e.Use(middlewarex.AuthMiddleware())

	fileHandler := file.NewHandler()
	// file apis
	fileGrp := e.Group("/file")
	fileGrp.GET("/:code/meta", fileHandler.GetMetaData)
	fileGrp.GET("/:code", fileHandler.Download)
	fileGrp.POST("/", fileHandler.Upload)
	fileGrp.DELETE("/:code", fileHandler.Delete)
}
