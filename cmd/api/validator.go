package main

import (
	"reflect"
	"strings"

	Config "github.com/Adhiana46/aegis-asesmen/config"
	PkgHttpServer "github.com/Adhiana46/aegis-asesmen/pkg/server/http"
	"github.com/go-playground/validator/v10"
)

func setupValidator(
	cfg *Config.Config,
	srv *PkgHttpServer.Server,
) {
	// Init Validator
	validatorInstance := validator.New()
	validatorInstance.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("name"), ",", 2)[0]
		if name == "" {
			name = strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		}
		if name == "-" {
			return ""
		}
		return name
	})

	srv.GetEngine().Validator = &echoValidator{
		validator: validatorInstance,
	}

}

type echoValidator struct {
	validator *validator.Validate
}

func (s *echoValidator) Validate(i interface{}) error {
	return s.validator.Struct(i)
}
