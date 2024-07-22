package model

import (
	"time"

	Entity "github.com/Adhiana46/aegis-asesmen/internal/user/domain/entity"
)

type User struct {
	Id        string    `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// turn user model into entity
func (m *User) ToEntity() *Entity.User {
	return &Entity.User{
		Id:        m.Id,
		Email:     m.Email,
		Password:  m.Password,
		Role:      m.Role,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// Create new user model from entity
func NewUserModel(e *Entity.User) *User {
	return &User{
		Id:        e.Id,
		Email:     e.Email,
		Password:  e.Password,
		Role:      e.Role,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}
