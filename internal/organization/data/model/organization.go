package model

import (
	"time"

	Entity "github.com/Adhiana46/aegis-asesmen/internal/organization/domain/entity"
)

type Organization struct {
	Id        string    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	CreatedBy string    `db:"created_by"`
	UpdatedAt time.Time `db:"updated_at"`
	UpdatedBy string    `db:"updated_by"`
}

// turn model into entity
func (m *Organization) ToEntity() *Entity.Organization {
	return &Entity.Organization{
		Id:        m.Id,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		CreatedBy: m.CreatedBy,
		UpdatedAt: m.UpdatedAt,
		UpdatedBy: m.UpdatedBy,
	}
}

// Create new model from entity
func NewOrganizationModel(e *Entity.Organization) *Organization {
	return &Organization{
		Id:        e.Id,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
		CreatedBy: e.CreatedBy,
		UpdatedAt: e.UpdatedAt,
		UpdatedBy: e.UpdatedBy,
	}
}
