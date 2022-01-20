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
	Id              int64   `json:"id"`
	Name            string  `json:"name"`
	Email           string  `json:"email"`
	Password        string  `json:"password"`
	Role            int64   `json:"role"`
	Status          int64   `json:"status"`
	InstanceId      *string `json:"instance_id"`
	InstanceName    *string `json:"instance_name"`
	DateCreated     string  `json:"date_created"`
	OldPassword     string  `json:"old_password"`
	DefaultPassword int64   `json:"default_password"`
}

type Filter struct {
	Role            string `json:"role"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Club            string `json:"id_instance"`
	Status          string `json:"status"`
	DefaultPassword string `json:"default_password"`
	SortBy          string `json:"sort_by"`
	SortDesc        string `json:"sort_desc"`
}

type ChangePassword struct {
	Id                   int64  `json:"user_id"`
	CurrentPassoword     string `json:"current_password"`
	NewPassword          string `json:"new_password"`
	ConfirmationPassword string `json:"confirmation_password"`
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
