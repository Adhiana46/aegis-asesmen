package event

import (
	"encoding/json"
	"fmt"

	UserEntity "github.com/Adhiana46/aegis-asesmen/internal/user/domain/entity"
	"github.com/pkg/errors"
)

// Topic
const (
	UserSigninTopic = "aegis.hiring"
)

type UserSigninEvent struct {
	UserId  string `json:"-"`
	Type    string `json:"type"`
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// GET TOPIC NAME
func (e *UserSigninEvent) Topic() string {
	return UserSigninTopic
}

// GET KEY NAME, used for partition
func (e *UserSigninEvent) Key() string {
	return e.UserId
}

// GET Payload
func (e *UserSigninEvent) Payload() ([]byte, error) {
	payload, err := json.Marshal(e)
	if err != nil {
		return nil, errors.Wrap(err, "json marshall failed")
	}

	return payload, nil
}

func NewUserSigninEvent(user *UserEntity.User) *UserSigninEvent {
	return &UserSigninEvent{
		UserId:  user.Id,
		Type:    "login",
		Status:  true,
		Message: fmt.Sprintf("%s logged in", user.Email),
	}
}
