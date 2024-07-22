package usecase

import (
	"context"
	"runtime"
	"strings"

	Config "github.com/Adhiana46/aegis-asesmen/config"
	Errors "github.com/Adhiana46/aegis-asesmen/errors"
	DTO "github.com/Adhiana46/aegis-asesmen/internal/organization/domain/dto"
	Repository "github.com/Adhiana46/aegis-asesmen/internal/organization/domain/repository"
	"github.com/pkg/errors"
)

type GetOrganizationByIdUsecase struct {
	config *Config.Config
	repo   Repository.IOrganizationRepository
}

func NewGetOrganizationByIdUsecase(
	config *Config.Config,
	repo Repository.IOrganizationRepository,
) *GetOrganizationByIdUsecase {
	return &GetOrganizationByIdUsecase{
		config: config,
		repo:   repo,
	}
}

func (u *GetOrganizationByIdUsecase) path() string {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	chunks := strings.Split(fn.Name(), ".")
	return "GetOrganizationByIdUsecase:" + chunks[len(chunks)-1]
}

func (u *GetOrganizationByIdUsecase) Do(ctx context.Context, id string) (*DTO.Organization, error) {
	path := u.path()

	obj, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, path)
	}
	if obj == nil {
		return nil, Errors.NewErrorDataNotFound()
	}

	output := DTO.NewOrganization(obj)

	return &output, nil
}
