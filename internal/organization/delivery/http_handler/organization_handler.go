package httphandler

import (
	"net/http"
	"strings"

	Config "github.com/Adhiana46/aegis-asesmen/config"
	Constants "github.com/Adhiana46/aegis-asesmen/constants"
	DTO "github.com/Adhiana46/aegis-asesmen/internal/organization/domain/dto"
	Usecase "github.com/Adhiana46/aegis-asesmen/internal/organization/domain/usecase"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type OrganizationHandler struct {
	cfg                        *Config.Config
	getListOrganizationUsecase *Usecase.GetListOrganizationUsecase
	getOrganizationByIdUsecase *Usecase.GetOrganizationByIdUsecase
	createOrganizationUsecase  *Usecase.CreateOrganizationUsecase
	updateOrganizationUsecase  *Usecase.UpdateOrganizationUsecase
	deleteOrganizationUsecase  *Usecase.DeleteOrganizationUsecase
}

func NewOrganizationHandler(
	cfg *Config.Config,
	getListOrganizationUsecase *Usecase.GetListOrganizationUsecase,
	getOrganizationByIdUsecase *Usecase.GetOrganizationByIdUsecase,
	createOrganizationUsecase *Usecase.CreateOrganizationUsecase,
	updateOrganizationUsecase *Usecase.UpdateOrganizationUsecase,
	deleteOrganizationUsecase *Usecase.DeleteOrganizationUsecase,
) *OrganizationHandler {
	return &OrganizationHandler{
		cfg:                        cfg,
		getListOrganizationUsecase: getListOrganizationUsecase,
		getOrganizationByIdUsecase: getOrganizationByIdUsecase,
		createOrganizationUsecase:  createOrganizationUsecase,
		updateOrganizationUsecase:  updateOrganizationUsecase,
		deleteOrganizationUsecase:  deleteOrganizationUsecase,
	}
}

func (h *OrganizationHandler) GetListOrganization(c echo.Context) error {
	bodyPayload := DTO.GetOrganizationListParam{}

	if err := c.Bind(&bodyPayload); err != nil {
		return errors.Wrap(err, "failed to parse body payload")
	}

	if err := c.Validate(&bodyPayload); err != nil {
		return errors.Wrap(err, "failed to validate body payload")
	}

	res, err := h.getListOrganizationUsecase.Do(c.Request().Context(), &bodyPayload)
	if err != nil {
		return errors.Wrap(err, "failed to get list organization")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status":  true,
		"message": Constants.MsgResponseOk,
		"data":    res,
	})
}

func (h *OrganizationHandler) GetOrganizationById(c echo.Context) error {
	resourceID := strings.TrimRight(c.Param("id"), "/")

	res, err := h.getOrganizationByIdUsecase.Do(c.Request().Context(), resourceID)
	if err != nil {
		return errors.Wrap(err, "failed to get organization by id")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status":  true,
		"message": Constants.MsgResponseOk,
		"data":    res,
	})
}

func (h *OrganizationHandler) CreateOrganization(c echo.Context) error {
	bodyPayload := DTO.CreateOrganizationParam{}

	if err := c.Bind(&bodyPayload); err != nil {
		return errors.Wrap(err, "failed to parse body payload")
	}

	if err := c.Validate(&bodyPayload); err != nil {
		return errors.Wrap(err, "failed to validate body payload")
	}

	err := h.createOrganizationUsecase.Do(c.Request().Context(), &bodyPayload)
	if err != nil {
		return errors.Wrap(err, "failed to create organization")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status":  true,
		"message": Constants.MsgResponseCreated,
		"data":    nil,
	})
}

func (h *OrganizationHandler) UpdateOrganization(c echo.Context) error {
	bodyPayload := DTO.UpdateOrganizationParam{}

	if err := c.Bind(&bodyPayload); err != nil {
		return errors.Wrap(err, "failed to parse body payload")
	}

	// trim slash
	bodyPayload.Id = strings.TrimRight(bodyPayload.Id, "/")

	if err := c.Validate(&bodyPayload); err != nil {
		return errors.Wrap(err, "failed to validate body payload")
	}

	err := h.updateOrganizationUsecase.Do(c.Request().Context(), &bodyPayload)
	if err != nil {
		return errors.Wrap(err, "failed to update organization")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status":  true,
		"message": Constants.MsgResponseOk,
		"data":    nil,
	})
}

func (h *OrganizationHandler) DeleteOrganization(c echo.Context) error {
	resourceID := strings.TrimRight(c.Param("id"), "/")

	err := h.deleteOrganizationUsecase.Do(c.Request().Context(), resourceID)
	if err != nil {
		return errors.Wrap(err, "delete drinking by id failed")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status":  true,
		"message": Constants.MsgResponseOk,
		"data":    nil,
	})
}
