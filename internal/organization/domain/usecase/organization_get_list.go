package usecase

import (
	"context"
	"math"
	"runtime"
	"strings"

	DTO "github.com/Adhiana46/aegis-asesmen/internal/organization/domain/dto"

	Config "github.com/Adhiana46/aegis-asesmen/config"
	Repository "github.com/Adhiana46/aegis-asesmen/internal/organization/domain/repository"
	"github.com/pkg/errors"
)

type GetListOrganizationUsecase struct {
	config *Config.Config
	repo   Repository.IOrganizationRepository
}

func NewGetListOrganizationUsecase(
	config *Config.Config,
	repo Repository.IOrganizationRepository,
) *GetListOrganizationUsecase {
	return &GetListOrganizationUsecase{
		config: config,
		repo:   repo,
	}
}

func (u *GetListOrganizationUsecase) path() string {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	chunks := strings.Split(fn.Name(), ".")
	return "GetListOrganizationUsecase:" + chunks[len(chunks)-1]
}

func (u *GetListOrganizationUsecase) Do(ctx context.Context, input *DTO.GetOrganizationListParam) (*DTO.GetOrganizationListResponse, error) {
	path := u.path()

	if input.Page == 0 {
		input.Page = 1
	}
	if input.Limit == 0 {
		input.Limit = 10
	}

	limit := input.Limit
	offset := (limit * input.Page) - limit

	entities, count, err := u.repo.GetList(ctx, offset, limit)
	if err != nil {
		return nil, errors.Wrap(err, path)
	}

	// entities -> DTOs
	dtos := make([]*DTO.Organization, len(entities))
	for i := 0; i < len(entities); i++ {
		d := DTO.NewOrganization(entities[i])
		dtos[i] = &d
	}

	// calc total page
	totalPage := 1
	if limit > 0 {
		totalPage = int(math.Ceil(float64(count) / float64(limit)))
	}

	output := DTO.GetOrganizationListResponse{
		CurrentPage: input.Page,
		TotalPage:   totalPage,
		TotalData:   count,
		Data:        dtos,
	}

	return &output, nil
}
