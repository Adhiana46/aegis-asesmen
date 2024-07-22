package dto

type UpdateOrganizationParam struct {
	Id   string `param:"id" name:"id" validate:"required"`
	Name string `json:"name" name:"name" validate:"required"`
}
