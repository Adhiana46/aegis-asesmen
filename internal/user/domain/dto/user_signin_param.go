package dto

type UserSigninParam struct {
	Email    string `json:"email" name:"email" validate:"required"`
	Password string `json:"password" name:"password" validate:"required"`
}
