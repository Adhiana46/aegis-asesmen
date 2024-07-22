package entity

import (
	"time"
)

type Organization struct {
	Id        string
	Name      string
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
	UpdatedBy string
}
