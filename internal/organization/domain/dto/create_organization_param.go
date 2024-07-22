package dto

type CreateOrganizationParam struct {
	Name string `json:"name" name:"name" validate:"required"`
}
