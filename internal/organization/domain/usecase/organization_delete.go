package usecase

import (
	"context"
	"runtime"
	"strings"

	Config "github.com/Adhiana46/aegis-asesmen/config"
	Constants "github.com/Adhiana46/aegis-asesmen/constants"
	Errors "github.com/Adhiana46/aegis-asesmen/errors"
	Repository "github.com/Adhiana46/aegis-asesmen/internal/organization/domain/repository"
	UserEntity "github.com/Adhiana46/aegis-asesmen/internal/user/domain/entity"
	"github.com/pkg/errors"
)

type DeleteOrganizationUsecase struct {
	config *Config.Config
	repo   Repository.IOrganizationRepository
}

func NewDeleteOrganizationUsecase(
	config *Config.Config,
	repo Repository.IOrganizationRepository,
) *DeleteOrganizationUsecase {
	return &DeleteOrganizationUsecase{
		config: config,
		repo:   repo,
	}
}

func (u *DeleteOrganizationUsecase) path() string {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	chunks := strings.Split(fn.Name(), ".")
	return "DeleteOrganizationUsecase:" + chunks[len(chunks)-1]
}

func (u *DeleteOrganizationUsecase) Do(ctx context.Context, id string) error {
	path := u.path()

	user, ok := ctx.Value("user").(*UserEntity.UserClaims)
	if !ok {
		return Errors.NewErrorInvalidCredentials()
	}

	obj, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, path)
	}
	if obj == nil {
		return Errors.NewErrorDataNotFound()
	}

	// Check Permission
	if user.Role != Constants.ROLE_SUPERADMIN && obj.CreatedBy != user.Id {
		return Errors.NewErrorInsufficientAccess()
	}

	if err := u.repo.Destroy(ctx, obj); err != nil {
		return errors.Wrap(err, path)
	}

	// SUCCESS
	return nil
}
