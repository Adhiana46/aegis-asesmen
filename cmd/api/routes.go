package main

import (
	"net/http"

	Config "github.com/Adhiana46/aegis-asesmen/config"
	UserHandler "github.com/Adhiana46/aegis-asesmen/internal/user/delivery/http_handler"
	PkgHttpServer "github.com/Adhiana46/aegis-asesmen/pkg/server/http"
	"github.com/labstack/echo/v4"
)

func setupHttpRoutes(
	cfg *Config.Config,
	srv *PkgHttpServer.Server,
	// Handlers
	userHandler *UserHandler.UserHandler,
) {
	// Set Routes
	err := srv.SetHandlers(func(e *echo.Echo) error {
		e.GET("/", func(c echo.Context) error {
			return c.String(http.StatusOK, cfg.App.Name+" "+cfg.App.Version)
		})

		authRouter := e.Group("/api/auth")
		authRouter.POST("/signin", userHandler.UserSignIn)

		return nil
	})
	if err != nil {
		panic(err)
	}
}
