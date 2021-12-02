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
	GetUsers(int, int, *users.Filter) (users.Users, *int, *rest_errors.RestErr)
	ChangePassword(users.ChangePassword) *rest_errors.RestErr
	CreateUser(users.User) (*users.User, *rest_errors.RestErr)
	UpdateUser(users.User) (*users.User, *rest_errors.RestErr)
	UpdateUserEdit(users.User) (*users.User, *rest_errors.RestErr)
	DeleteUser(user users.User) *rest_errors.RestErr
	GetAttendants(search string) (users.Users, *rest_errors.RestErr)
}

func (s *usersService) GetUser(userId int64) (*users.User, *rest_errors.RestErr) {
	result := &users.User{Id: userId}
	if err := result.GetUser(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *usersService) GetUsers(page int, itemsPerPage int, filter *users.Filter) (users.Users, *int, *rest_errors.RestErr) {
	dao := &users.User{}
	users, total, err := dao.GetUsers(page, itemsPerPage, filter)
	if err != nil {
		return nil, nil, err
	}

	return users, total, nil
}

func (s *usersService) GetAttendants(search string) (users.Users, *rest_errors.RestErr) {
	dao := &users.User{}
	users, err := dao.GetAttendants(search)
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
	current.InstanceId = user.InstanceId
	current.Name = user.Name

	if err := current.Update(); err != nil {
		return nil, err
	}

	return current, nil
}

func (s *usersService) UpdateUserEdit(user users.User) (*users.User, *rest_errors.RestErr) {
	current, err := UsersService.GetUser(user.Id)
	if err != nil {
		return nil, err
	}

	if user.Name != "" {
		current.Name = user.Name

		if err := current.UpdateUserName(); err != nil {
			return nil, err
		}
	}

	if user.Email != "" {
		current.Email = user.Email

		if err := current.UpdateUserEmail(); err != nil {
			return nil, err
		}
	}

	if current.Password == crypto_utils.GetMd5(user.Password) {
		return nil, rest_errors.NewBadRequestError("A nova senha deve ser diferente da senha atual")
	} else if current.Password == crypto_utils.GetMd5(user.OldPassword) {
		current.Password = crypto_utils.GetMd5(user.Password)
		if err := current.ChangePassword(); err != nil {
			return nil, err
		}
	} else if user.Password != "" {
		return nil, rest_errors.NewBadRequestError("A senha antiga não corresponde")
	}

	return current, nil
}

func (s *usersService) ChangePassword(user users.ChangePassword) *rest_errors.RestErr {
	current, err := UsersService.GetUser(user.Id)
	if err != nil {
		return rest_errors.NewBadRequestError("Usuário não encontrado")
	}

	if user.CurrentPassoword != user.CurrentPassoword {
		return rest_errors.NewBadRequestError("Senha atual inválida")
	}

	if user.ConfirmationPassword != user.NewPassword {
		return rest_errors.NewBadRequestError("Senha/senha de confirmação devem ser iguais")
	}

	if len(user.NewPassword) < 8 {
		return rest_errors.NewBadRequestError("A senha deve possuir no minimo 8 caracteres")
	}

	current.Password = crypto_utils.GetMd5(user.NewPassword)

	if err := current.ChangePassword(); err != nil {
		return nil
	}

	return nil
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
