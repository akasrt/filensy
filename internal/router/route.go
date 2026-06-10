package router

import (
	"github.com/akasrt/filensy/internal/file"
	"github.com/akasrt/filensy/internal/middlewarex"
	"github.com/labstack/echo/v5"
)

func AddRoutes(e *echo.Echo) {
	e.Use(middlewarex.AuthMiddleware())
	apiGrp := e.Group("/api/v1")

	fileHandler := file.NewHandler()
	// file apis
	fileGrp := apiGrp.Group("/file")
	fileGrp.GET("/:code/meta", fileHandler.GetMetaData)
	fileGrp.GET("/:code", fileHandler.Download)
	fileGrp.POST("", fileHandler.Upload)
	fileGrp.DELETE("/:code", fileHandler.Delete)
}
