package main

import (
	Config "github.com/Adhiana46/aegis-asesmen/config"
	PkgHttpServer "github.com/Adhiana46/aegis-asesmen/pkg/server/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func setupMiddlewares(
	cfg *Config.Config,
	srv *PkgHttpServer.Server,
) {
	srv.GetEngine().Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
	}))
	srv.GetEngine().Use(middleware.Logger())
}
