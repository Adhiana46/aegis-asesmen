package source

import (
	"context"

	Model "github.com/Adhiana46/aegis-asesmen/internal/user/data/model"
)

type IUserPersistent interface {
	GetByID(ctx context.Context, id string) (*Model.User, error)
	GetByEmail(ctx context.Context, email string) (*Model.User, error)
	Store(ctx context.Context, userModel *Model.User) error
	Update(ctx context.Context, userModel *Model.User) error
	Destroy(ctx context.Context, userModel *Model.User) error
}
