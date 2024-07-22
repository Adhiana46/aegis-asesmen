package usecase

import (
	"context"
	"runtime"
	"strings"
	"time"

	Config "github.com/Adhiana46/aegis-asesmen/config"
	Errors "github.com/Adhiana46/aegis-asesmen/errors"
	DTO "github.com/Adhiana46/aegis-asesmen/internal/organization/domain/dto"
	Entity "github.com/Adhiana46/aegis-asesmen/internal/organization/domain/entity"
	Repository "github.com/Adhiana46/aegis-asesmen/internal/organization/domain/repository"
	"github.com/oklog/ulid/v2"
	"github.com/pkg/errors"
)

type CreateOrganizationUsecase struct {
	config *Config.Config
	repo   Repository.IOrganizationRepository
}

func NewCreateOrganizationUsecase(
	config *Config.Config,
	repo Repository.IOrganizationRepository,
) *CreateOrganizationUsecase {
	return &CreateOrganizationUsecase{
		config: config,
		repo:   repo,
	}
}

func (u *CreateOrganizationUsecase) path() string {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	chunks := strings.Split(fn.Name(), ".")
	return "CreateOrganizationUsecase:" + chunks[len(chunks)-1]
}

func (u *CreateOrganizationUsecase) Do(ctx context.Context, input *DTO.CreateOrganizationParam) error {
	path := u.path()

	existingObj, err := u.repo.GetByName(ctx, input.Name)
	if err != nil {
		return errors.Wrap(err, path)
	}
	if existingObj != nil {
		return Errors.NewErrorDataAlreadyExists()
	}

	obj := Entity.Organization{
		Id:        ulid.Make().String(),
		Name:      input.Name,
		CreatedAt: time.Now(),
		CreatedBy: "", // TODO: user_id
		UpdatedAt: time.Now(),
		UpdatedBy: "", // TODO: user_id
	}

	// save to repo
	if err := u.repo.Store(ctx, &obj); err != nil {
		return errors.Wrap(err, path)
	}

	return nil
}
