package repository

import (
	"context"
	"runtime"
	"strings"

	Model "github.com/Adhiana46/aegis-asesmen/internal/organization/data/model"
	DataSource "github.com/Adhiana46/aegis-asesmen/internal/organization/data/source"
	Entity "github.com/Adhiana46/aegis-asesmen/internal/organization/domain/entity"
	"github.com/pkg/errors"
)

type IOrganizationRepository interface {
	GetList(ctx context.Context, offset, limit int) ([]*Entity.Organization, int, error)
	GetByID(ctx context.Context, id string) (*Entity.Organization, error)
	GetByName(ctx context.Context, name string) (*Entity.Organization, error)
	Store(ctx context.Context, entity *Entity.Organization) error
	Update(ctx context.Context, entity *Entity.Organization) error
	Destroy(ctx context.Context, entity *Entity.Organization) error
}

type organizationRepository struct {
	persistent DataSource.IOrganizationPersistent
}

func NewOrganizationRepository(
	persistent DataSource.IOrganizationPersistent,
) IOrganizationRepository {
	return &organizationRepository{
		persistent: persistent,
	}
}

func (u *organizationRepository) path() string {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	chunks := strings.Split(fn.Name(), ".")
	return "organizationRepository:" + chunks[len(chunks)-1]
}

func (r *organizationRepository) GetList(ctx context.Context, offset, limit int) ([]*Entity.Organization, int, error) {
	path := r.path()

	rows, err := r.persistent.GetList(ctx, offset, limit)
	if err != nil {
		return nil, 0, errors.Wrap(err, path)
	}

	totalRows, err := r.persistent.CountList(ctx)
	if err != nil {
		return nil, 0, errors.Wrap(err, path)
	}

	entities := make([]*Entity.Organization, len(rows))
	for i := 0; i < len(rows); i++ {
		entities[i] = rows[i].ToEntity()
	}

	return entities, totalRows, nil

}

func (r *organizationRepository) GetByID(ctx context.Context, id string) (*Entity.Organization, error) {
	path := r.path()

	// GET FROM PERSISTENT
	model, err := r.persistent.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, path)
	}
	if model == nil {
		return nil, nil
	}

	entity := model.ToEntity()

	return entity, nil
}

func (r *organizationRepository) GetByName(ctx context.Context, name string) (*Entity.Organization, error) {
	path := r.path()

	// GET FROM PERSISTENT
	model, err := r.persistent.GetByName(ctx, name)
	if err != nil {
		return nil, errors.Wrap(err, path)
	}
	if model == nil {
		return nil, nil
	}

	entity := model.ToEntity()

	return entity, nil
}

func (r *organizationRepository) Store(ctx context.Context, entity *Entity.Organization) error {
	path := r.path()

	model := Model.NewOrganizationModel(entity)

	err := r.persistent.Store(ctx, model)
	if err != nil {
		return errors.Wrap(err, path)
	}

	return nil
}

func (r *organizationRepository) Update(ctx context.Context, entity *Entity.Organization) error {
	path := r.path()

	model := Model.NewOrganizationModel(entity)

	err := r.persistent.Update(ctx, model)
	if err != nil {
		return errors.Wrap(err, path)
	}

	return nil
}

func (r *organizationRepository) Destroy(ctx context.Context, entity *Entity.Organization) error {
	path := r.path()

	model := Model.NewOrganizationModel(entity)

	err := r.persistent.Destroy(ctx, model)
	if err != nil {
		return errors.Wrap(err, path)
	}

	return nil
}
