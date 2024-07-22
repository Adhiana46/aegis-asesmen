package main

import (
	"net/http"

	Config "github.com/Adhiana46/aegis-asesmen/config"
	OrganizationHandler "github.com/Adhiana46/aegis-asesmen/internal/organization/delivery/http_handler"
	UserHandler "github.com/Adhiana46/aegis-asesmen/internal/user/delivery/http_handler"
	PkgHttpServer "github.com/Adhiana46/aegis-asesmen/pkg/server/http"
	"github.com/labstack/echo/v4"
)

func setupHttpRoutes(
	cfg *Config.Config,
	srv *PkgHttpServer.Server,
	// Handlers
	userHandler *UserHandler.UserHandler,
	organizationHandler *OrganizationHandler.OrganizationHandler,
) {
	// Set Routes
	err := srv.SetHandlers(func(e *echo.Echo) error {
		e.GET("/", func(c echo.Context) error {
			return c.String(http.StatusOK, cfg.App.Name+" "+cfg.App.Version)
		})

		authRouter := e.Group("/api/auth")
		authRouter.POST("/signin", userHandler.UserSignIn)

		organizationRouter := e.Group("/api/organization")
		organizationRouter.GET("/", organizationHandler.GetListOrganization)
		organizationRouter.GET("/:id", organizationHandler.GetOrganizationById)
		organizationRouter.POST("/", organizationHandler.CreateOrganization)
		organizationRouter.PUT("/:id", organizationHandler.UpdateOrganization)
		organizationRouter.DELETE("/:id", organizationHandler.DeleteOrganization)

		return nil
	})
	if err != nil {
		panic(err)
	}
}
