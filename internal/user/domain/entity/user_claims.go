package entity

import (
	"github.com/golang-jwt/jwt/v5"
)

// JWT TOKEN PAYLOAD
type UserClaims struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}
