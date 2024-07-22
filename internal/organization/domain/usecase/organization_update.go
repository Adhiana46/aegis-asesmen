package usecase

import (
	"context"
	"runtime"
	"strings"
	"time"

	Config "github.com/Adhiana46/aegis-asesmen/config"
	Errors "github.com/Adhiana46/aegis-asesmen/errors"
	DTO "github.com/Adhiana46/aegis-asesmen/internal/organization/domain/dto"
	Repository "github.com/Adhiana46/aegis-asesmen/internal/organization/domain/repository"
	"github.com/pkg/errors"
)

type UpdateOrganizationUsecase struct {
	config *Config.Config
	repo   Repository.IOrganizationRepository
}

func NewUpdateOrganizationUsecase(
	config *Config.Config,
	repo Repository.IOrganizationRepository,
) *UpdateOrganizationUsecase {
	return &UpdateOrganizationUsecase{
		config: config,
		repo:   repo,
	}
}

func (u *UpdateOrganizationUsecase) path() string {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	chunks := strings.Split(fn.Name(), ".")
	return "UpdateOrganizationUsecase:" + chunks[len(chunks)-1]
}

func (u *UpdateOrganizationUsecase) Do(ctx context.Context, input *DTO.UpdateOrganizationParam) error {
	path := u.path()

	obj, err := u.repo.GetByID(ctx, input.Id)
	if err != nil {
		return errors.Wrap(err, path)
	}
	if obj == nil {
		return Errors.NewErrorDataNotFound()
	}

	// check by name
	objByName, err := u.repo.GetByName(ctx, input.Name)
	if err != nil {
		return errors.Wrap(err, path)
	}
	if objByName != nil && objByName.Id != input.Id {
		return Errors.NewErrorDataAlreadyExists()
	}

	// Update Values
	obj.Name = input.Name
	obj.UpdatedAt = time.Now()
	obj.UpdatedBy = "" // TODO: user_id

	if err := u.repo.Update(ctx, obj); err != nil {
		return errors.Wrap(err, path)
	}

	// SUCCESS
	return nil
}
