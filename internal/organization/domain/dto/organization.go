package dto

import (
	Entity "github.com/Adhiana46/aegis-asesmen/internal/organization/domain/entity"
)

type Organization struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func NewOrganization(entity *Entity.Organization) Organization {
	return Organization{
		Id:   entity.Id,
		Name: entity.Name,
	}
}
