package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	Config "github.com/Adhiana46/aegis-asesmen/config"
	Errors "github.com/Adhiana46/aegis-asesmen/errors"
	OrganizationHandler "github.com/Adhiana46/aegis-asesmen/internal/organization/delivery/http_handler"
	UserHandler "github.com/Adhiana46/aegis-asesmen/internal/user/delivery/http_handler"
	UserEntity "github.com/Adhiana46/aegis-asesmen/internal/user/domain/entity"
	PkgHttpServer "github.com/Adhiana46/aegis-asesmen/pkg/server/http"
	"github.com/golang-jwt/jwt/v5"
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

		organizationRouter := e.Group("/api/organization", authMiddleware(cfg))
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

func authMiddleware(cfg *Config.Config) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return Errors.NewErrorInvalidToken()
			}

			authHeaderChunk := strings.Split(authHeader, " ")
			if len(authHeaderChunk) != 2 {
				return Errors.NewErrorInvalidToken()
			}

			tokenType, tokenStr := authHeaderChunk[0], authHeaderChunk[1]
			if strings.ToLower(tokenType) != "bearer" {
				return Errors.NewErrorInvalidToken()
			}

			token, claims, err := parseToken(cfg, tokenStr)
			if err != nil || !token.Valid {
				return Errors.NewErrorInvalidToken()
			}

			if claims.RegisteredClaims.Subject != "access-token" {
				return Errors.NewErrorInvalidToken()
			}

			// Set user
			c.Set("user", claims)

			// Set request context
			ctx := context.WithValue(c.Request().Context(), "user", claims)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

func parseToken(cfg *Config.Config, tokenStr string) (*jwt.Token, *UserEntity.UserClaims, error) {
	secretKey := cfg.JWT.SecretKey

	// VALIDATE tokenStr Format
	_, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return nil, nil, Errors.NewErrorInvalidToken()
	}

	token, err := jwt.ParseWithClaims(tokenStr, &UserEntity.UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		// validate signing algo
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, nil, err
	}

	if !token.Valid {
		return nil, nil, Errors.NewErrorInvalidToken()
	}

	claims, ok := token.Claims.(*UserEntity.UserClaims)
	if !ok {
		return nil, nil, Errors.NewErrorInvalidToken()
	}

	return token, claims, nil
}
