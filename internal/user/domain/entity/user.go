package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (u *User) IsPasswordMatch(plain string) bool {
	byteHash := []byte(u.Password)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(plain))

	return err == nil
}

func (u *User) SetPassword(plain string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.MinCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}
