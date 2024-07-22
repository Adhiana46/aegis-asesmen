package repository

import (
	"context"
	"runtime"
	"strings"

	Model "github.com/Adhiana46/aegis-asesmen/internal/user/data/model"
	DataSource "github.com/Adhiana46/aegis-asesmen/internal/user/data/source"
	Entity "github.com/Adhiana46/aegis-asesmen/internal/user/domain/entity"
	"github.com/pkg/errors"
)

type IUserRepository interface {
	GetByID(ctx context.Context, id string) (*Entity.User, error)
	GetByEmail(ctx context.Context, email string) (*Entity.User, error)
	Store(ctx context.Context, user *Entity.User) error
	Update(ctx context.Context, user *Entity.User) error
	Destroy(ctx context.Context, user *Entity.User) error
}

type userRepository struct {
	persistent DataSource.IUserPersistent
}

func NewUserRepository(
	persistent DataSource.IUserPersistent,
) IUserRepository {
	return &userRepository{
		persistent: persistent,
	}
}

func (u *userRepository) path() string {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	chunks := strings.Split(fn.Name(), ".")
	return "userRepository:" + chunks[len(chunks)-1]
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*Entity.User, error) {
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

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*Entity.User, error) {
	path := r.path()

	// GET FROM PERSISTENT
	model, err := r.persistent.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, path)
	}
	if model == nil {
		return nil, nil
	}

	entity := model.ToEntity()

	return entity, nil
}

func (r *userRepository) Store(ctx context.Context, user *Entity.User) error {
	path := r.path()

	model := Model.NewUserModel(user)

	err := r.persistent.Store(ctx, model)
	if err != nil {
		return errors.Wrap(err, path)
	}

	return nil
}

func (r *userRepository) Update(ctx context.Context, user *Entity.User) error {
	path := r.path()

	model := Model.NewUserModel(user)

	err := r.persistent.Update(ctx, model)
	if err != nil {
		return errors.Wrap(err, path)
	}

	return nil
}

func (r *userRepository) Destroy(ctx context.Context, user *Entity.User) error {
	path := r.path()

	model := Model.NewUserModel(user)

	err := r.persistent.Destroy(ctx, model)
	if err != nil {
		return errors.Wrap(err, path)
	}

	return nil
}
