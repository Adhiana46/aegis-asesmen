package usecase

import (
	"context"
	"runtime"
	"strings"
	"time"

	Config "github.com/Adhiana46/aegis-asesmen/config"
	Errors "github.com/Adhiana46/aegis-asesmen/errors"
	DTO "github.com/Adhiana46/aegis-asesmen/internal/user/domain/dto"
	Entity "github.com/Adhiana46/aegis-asesmen/internal/user/domain/entity"
	Event "github.com/Adhiana46/aegis-asesmen/internal/user/domain/event"
	Repository "github.com/Adhiana46/aegis-asesmen/internal/user/domain/repository"
	PkgKafkaPublisher "github.com/Adhiana46/aegis-asesmen/pkg/kafka/publisher"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type UserSigninUsecase struct {
	config    *Config.Config
	publisher *PkgKafkaPublisher.Publisher
	userRepo  Repository.IUserRepository
}

func NewUserSigninUsecase(
	config *Config.Config,
	publisher *PkgKafkaPublisher.Publisher,
	userRepo Repository.IUserRepository,
) *UserSigninUsecase {
	return &UserSigninUsecase{
		config:    config,
		publisher: publisher,
		userRepo:  userRepo,
	}
}

func (u *UserSigninUsecase) path() string {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	chunks := strings.Split(fn.Name(), ".")
	return "UserSigninUsecase:" + chunks[len(chunks)-1]
}

func (u *UserSigninUsecase) Do(ctx context.Context, input *DTO.UserSigninParam) (*DTO.UserTokenResponse, error) {
	path := u.path()

	// FIND USER BY EMAIL
	user, err := u.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.Wrap(err, path)
	}
	if user == nil {
		return nil, Errors.NewErrorInvalidCredentials()
	}

	if !user.IsPasswordMatch(input.Password) {
		return nil, Errors.NewErrorInvalidCredentials()
	}

	// GENERATE TOKEN
	accessToken, err := u.generateToken(user)
	if err != nil {
		return nil, errors.Wrap(err, path)
	}

	output := DTO.UserTokenResponse{
		AccessToken: accessToken,
	}

	//  Publish Event
	evt := Event.NewUserSigninEvent(user)
	if err := u.publisher.Publish(evt); err != nil {
		return nil, errors.Wrap(err, path)
	}

	// SUCCESS
	return &output, nil
}

// Generate access & refresh token
func (u *UserSigninUsecase) generateToken(user *Entity.User) (string, error) {
	path := u.path()

	claims := Entity.UserClaims{
		Id:    user.Id,
		Email: user.Email,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: "access-token",
			Issuer:  u.config.JWT.Issuer,
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(time.Hour), // 1 hour
			},
			IssuedAt: &jwt.NumericDate{
				Time: time.Now(),
			},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(u.config.JWT.SecretKey))
	if err != nil {
		return "", errors.Wrap(err, path)
	}

	return tokenStr, nil
}
