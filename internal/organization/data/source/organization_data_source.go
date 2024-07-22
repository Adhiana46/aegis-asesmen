package source

import (
	"context"

	Model "github.com/Adhiana46/aegis-asesmen/internal/organization/data/model"
)

type IOrganizationPersistent interface {
	GetList(ctx context.Context, offset, limit int) ([]*Model.Organization, error)
	CountList(ctx context.Context) (int, error)
	GetByID(ctx context.Context, id string) (*Model.Organization, error)
	GetByName(ctx context.Context, name string) (*Model.Organization, error)
	Store(ctx context.Context, model *Model.Organization) error
	Update(ctx context.Context, model *Model.Organization) error
	Destroy(ctx context.Context, model *Model.Organization) error
}
