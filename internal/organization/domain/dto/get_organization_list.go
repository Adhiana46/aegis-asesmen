package dto

type GetOrganizationListParam struct {
	Page  int `query:"page" name:"page" validate:""`
	Limit int `query:"limit" name:"limit" validate:""`
}

type GetOrganizationListResponse struct {
	CurrentPage int             `json:"current_page"`
	TotalPage   int             `json:"total_page"`
	TotalData   int             `json:"total_data"`
	Data        []*Organization `json:"data"`
}
