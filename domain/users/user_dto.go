package users

import (
	"strings"

	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
)

const (
	StatusActive = 1
)

type Users []User

type User struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Role        int64  `json:"role"`
	Status      int64  `json:"status"`
	DateCreated string `json:"date_created"`
}

func (user *User) Validate() *rest_errors.RestErr {
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" {
		return rest_errors.NewBadRequestError("Invalid email address")
	}

	if user.Email == "" {
		return rest_errors.NewBadRequestError("Invalid password")
	}

	return nil
}
