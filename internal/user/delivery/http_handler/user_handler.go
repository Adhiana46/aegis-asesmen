package httphandler

import (
	"net/http"

	Config "github.com/Adhiana46/aegis-asesmen/config"
	Constants "github.com/Adhiana46/aegis-asesmen/constants"
	DTO "github.com/Adhiana46/aegis-asesmen/internal/user/domain/dto"
	Usecase "github.com/Adhiana46/aegis-asesmen/internal/user/domain/usecase"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type UserHandler struct {
	cfg               *Config.Config
	userSigninUsecase *Usecase.UserSigninUsecase
}

func NewAuthHandler(
	cfg *Config.Config,
	userSigninUsecase *Usecase.UserSigninUsecase,
) *UserHandler {
	return &UserHandler{
		cfg:               cfg,
		userSigninUsecase: userSigninUsecase,
	}
}

func (h *UserHandler) UserSignIn(c echo.Context) error {
	bodyPayload := DTO.UserSigninParam{}

	if err := c.Bind(&bodyPayload); err != nil {
		return errors.Wrap(err, "failed to parse body payload")
	}

	if err := c.Validate(&bodyPayload); err != nil {
		return errors.Wrap(err, "failed to validate body payload")
	}

	res, err := h.userSigninUsecase.Do(c.Request().Context(), &bodyPayload)
	if err != nil {
		return errors.Wrap(err, "user signin failed")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status":  true,
		"message": Constants.MsgSuccessSignin,
		"data":    res,
	})
}
