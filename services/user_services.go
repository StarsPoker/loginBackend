package services

import (
	"github.com/StarsPoker/loginBackend/domain/users"
	"github.com/StarsPoker/loginBackend/utils/crypto_utils.go"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
	//"net/http"
)

var (
	UsersService usersServiceInterface = &usersService{}
)

type usersService struct {
}

type usersServiceInterface interface {
	GetUser(int64) (*users.User, *rest_errors.RestErr)
	GetUsers(int, int) (users.Users, *int, *rest_errors.RestErr)
	CreateUser(users.User) (*users.User, *rest_errors.RestErr)
	UpdateUser(users.User) (*users.User, *rest_errors.RestErr)
	DeleteUser(user users.User) *rest_errors.RestErr
	GetAttendances(search string) (users.Users, *rest_errors.RestErr)
}

func (s *usersService) GetUser(userId int64) (*users.User, *rest_errors.RestErr) {
	result := &users.User{Id: userId}
	if err := result.GetUser(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *usersService) GetUsers(page int, itemsPerPage int) (users.Users, *int, *rest_errors.RestErr) {
	dao := &users.User{}
	users, total, err := dao.GetUsers(page, itemsPerPage)
	if err != nil {
		return nil, nil, err
	}

	return users, total, nil
}

func (s *usersService) GetAttendances(search string) (users.Users, *rest_errors.RestErr) {
	dao := &users.User{}
	users, err := dao.GetAttendances(search)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *usersService) CreateUser(user users.User) (*users.User, *rest_errors.RestErr) {

	if err := user.Validate(); err != nil {
		return nil, err
	}

	user.Password = crypto_utils.GetMd5("123456")

	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *usersService) UpdateUser(user users.User) (*users.User, *rest_errors.RestErr) {
	current, err := UsersService.GetUser(user.Id)
	if err != nil {
		return nil, err
	}

	current.Email = user.Email
	current.Role = user.Role
	current.Status = user.Status

	if err := current.Update(); err != nil {
		return nil, err
	}

	return current, nil
}

func (s *usersService) DeleteUser(user users.User) *rest_errors.RestErr {
	current, err := UsersService.GetUser(user.Id)
	if err != nil {
		return err
	}

	if err := current.Delete(); err != nil {
		return nil
	}

	return nil
}
