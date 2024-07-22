package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	Config "github.com/Adhiana46/aegis-asesmen/config"
	Errors "github.com/Adhiana46/aegis-asesmen/errors"
	PkgHttpServer "github.com/Adhiana46/aegis-asesmen/pkg/server/http"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func setupErrorHandler(
	cfg *Config.Config,
	srv *PkgHttpServer.Server,
) {
	srv.GetEngine().HTTPErrorHandler = func(err error, c echo.Context) {
		var (
			statusCode int
			message    string
			errorsData any
		)
		stackTraces := []string{
			err.Error(),
		}
		stackTraces = append(stackTraces, strings.Split(strings.ReplaceAll(fmt.Sprintf("%+v", err), "\t", "  "), "\n")...)

		// Handle Errors
		switch {
		case errors.As(err, &validator.ValidationErrors{}):
			errValidation := validator.ValidationErrors{}
			errors.As(err, &errValidation)

			message = "Validation Error"
			statusCode = 400
			errorsData = func(validationErrs validator.ValidationErrors) map[string][]string {
				errorFields := map[string][]string{}
				for _, e := range validationErrs {
					errorFields[e.Field()] = append(errorFields[e.Field()], e.Tag())
				}

				return errorFields
			}(errValidation)
		case errors.As(err, new(Errors.InternalError)):
			errInternal := Errors.NewInternalError()
			errors.As(err, &errInternal)

			message = errInternal.Error()
			statusCode = errInternal.HttpStatusCode()
		default:
			message = "Internal Server Error"
			statusCode = http.StatusInternalServerError
		}

		details := map[string]any{
			"method":        c.Request().Method,
			"endpoint":      c.Request().URL.String(),
			"client_ip":     c.RealIP(),
			"user_agent":    c.Request().UserAgent(),
			"payload":       c.Get("request_body"),
			"error_message": err.Error(),
			"stack_trace":   stackTraces,
		}
		slog.Error("Error processing HTTP request", slog.Any("details", details))

		err = c.JSON(statusCode, map[string]any{
			"status":       false,
			"message":      message,
			"errors":       errorsData,
			"stack_traces": stackTraces,
		})
		if err != nil {
			slog.Error("Error writing response json", slog.String("error", err.Error()))
		}
	}
}
