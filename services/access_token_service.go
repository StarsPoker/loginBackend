package services

import (
	"fmt"
	"strings"

	"github.com/StarsPoker/loginBackend/domain/access_token"
	"github.com/StarsPoker/loginBackend/domain/one_time_password"
	"github.com/StarsPoker/loginBackend/domain/users"

	"github.com/StarsPoker/loginBackend/domain/chat_repository"
	"github.com/StarsPoker/loginBackend/utils/crypto_utils"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
)

var (
	AccessTokenService AccessTokenServiceInterface = &accessTokenService{}
)

type accessTokenService struct {
}

type AccessTokenServiceInterface interface {
	GetById(string) (*access_token.AccessToken, *rest_errors.RestErr)
	Create(accessTokenRequest access_token.AccessTokenRequest) (*one_time_password.OneTimePassword, *rest_errors.RestErr)
	ValidateAccessToken(string) *rest_errors.RestErr
	Delete(string) *rest_errors.RestErr
	CheckAuth(accessTokenRequest access_token.AccessTokenRequest, host string, client_ip string) (*one_time_password.OneTimePassword, *rest_errors.RestErr)
	CreateDevelopment(accessTokenRequest access_token.AccessTokenRequest, host string, client_ip string) (*one_time_password.OneTimePassword, *rest_errors.RestErr)
	DeleteExpiredAccesTokens()
	DeleteExpiredOneTimePasswords()
}

func (s *accessTokenService) GetById(accessTokenId string) (*access_token.AccessToken, *rest_errors.RestErr) {
	if len(accessTokenId) == 0 {
		return nil, rest_errors.NewBadRequestError("invalid access token")
	}
	if strings.Contains(accessTokenId, "Bearer") {
		accessTokenId = accessTokenId[7:]
	}
	accessTokenId = strings.TrimSpace(accessTokenId)

	tkn, errGet := access_token.CheckToken(accessTokenId)
	if errGet != nil {
		return nil, errGet
	}
	// tkn, err := access_token.GetById(accessTokenId)
	// if err != nil {
	// 	return nil, err
	// }
	return tkn, nil
}

func (s *accessTokenService) Create(accessTokenRequest access_token.AccessTokenRequest) (*one_time_password.OneTimePassword, *rest_errors.RestErr) {
	if err := accessTokenRequest.Validate(); err != nil {
		return nil, err
	}

	user := &users.User{
		Email:  accessTokenRequest.Username,
		Status: users.StatusActive,
	}

	if err := user.FindByEmailAndPassword(); err != nil {
		return nil, err
	}

	if crypto_utils.GetMd5(accessTokenRequest.Password) != user.Password {
		return nil, rest_errors.NewInternalServerError("invalid credentials")
	}

	errGetUser := user.GetUser()
	if errGetUser != nil {
		return nil, errGetUser
	}

	if user.Role == nil {
		return nil, rest_errors.NewInternalServerError("Usuário não possui perfil de acesso associado")
	}

	var otp one_time_password.OneTimePassword
	otp = otp.CreateOtp(user)

	err := one_time_password.Insert(otp)

	if err != nil {
		return nil, err
	}

	/* if user.Contact != nil {
		go chat_repository.SendWhatsappMessage(otp, user)
	} */
	chat_repository.SendMail(otp, user)

	otp.Code = "anonimized"
	return &otp, nil
}

func (s *accessTokenService) ValidateAccessToken(accessTokenId string) *rest_errors.RestErr {
	accessTokenId = strings.Replace(accessTokenId, "Bearer ", "", 1)
	if accessTokenId == "" {
		fmt.Println("access token not found")
		return rest_errors.NewUnauthorizedError("access token not found")
	}
	accessToken, err := AccessTokenService.GetById(accessTokenId)
	if err != nil {
		fmt.Println("invalid access token")
		return rest_errors.NewUnauthorizedError("invalid access token")
	}

	var user users.User
	if accessToken.UserHost == "localhost" || strings.Contains(accessToken.UserHost, "192.168.1") || strings.Contains(accessToken.UserHost, "192.168.2") {
		if err := user.ValidateExternalAccess(accessToken.UserId); err != nil {
			return rest_errors.NewUnauthorizedError("blocked external access")
		}
	}

	expired := accessToken.IsExpired()
	if expired {
		_ = access_token.Delete(accessTokenId)
		fmt.Println("access token expired")
		return rest_errors.NewUnauthorizedError("access token expired")
	}

	_ = access_token.UpdateLastInteraction(accessTokenId)
	fmt.Println("access token valid")
	return nil
}

func (s *accessTokenService) Delete(accessTokenId string) *rest_errors.RestErr {
	err := access_token.Delete(accessTokenId)
	if err != nil {
		return err
	}

	return nil
}

func (s *accessTokenService) CheckAuth(accessTokenRequest access_token.AccessTokenRequest, host string, client_ip string) (*one_time_password.OneTimePassword, *rest_errors.RestErr) {

	otp, err := one_time_password.GetAuth(accessTokenRequest)
	if err != nil {
		return nil, err
	}

	otp, err = one_time_password.CheckAuth(otp, accessTokenRequest)
	if err != nil {
		return nil, err
	}

	user := &users.User{
		Id: otp.UserId,
	}

	errGetUser := user.GetUser()
	if errGetUser != nil {
		return nil, errGetUser
	}

	rec := false //recursive calling
	rec = one_time_password.DeleteAllById(*otp, rec)
	if !rec {
		return nil, nil
	}

	at := access_token.GetNewAccessToken(user.Id, *user.Role)
	at.Generate()
	at.UserHost = host
	at.UserClientIp = client_ip
	at.UserIpFront = accessTokenRequest.UserIpFront
	err = access_token.Create(at)
	if err != nil {
		return nil, err
	}
	otp.AccessToken = at

	return otp, nil
}
func (s *accessTokenService) CreateDevelopment(accessTokenRequest access_token.AccessTokenRequest, host string, client_ip string) (*one_time_password.OneTimePassword, *rest_errors.RestErr) {
	otp := &one_time_password.OneTimePassword{
		AccessToken: access_token.AccessToken{
			AccessToken: "",
		},
	}
	user := &users.User{
		Email:  accessTokenRequest.Username,
		Status: users.StatusActive,
	}

	if err := user.FindByEmailAndPassword(); err != nil {
		return nil, err
	}

	at := access_token.GetNewAccessToken(user.Id, *user.Role)
	at.Generate()
	at.UserHost = host
	at.UserClientIp = client_ip
	at.UserIpFront = accessTokenRequest.UserIpFront
	err := access_token.Create(at)

	if err != nil {
		return nil, err
	}

	otp.AccessToken = at
	return otp, nil
}

func (s *accessTokenService) DeleteExpiredAccesTokens() {
	err := access_token.DeleteExpiredAccesTokens()
	if err != nil {
		fmt.Println(err)
	}
}

func (s *accessTokenService) DeleteExpiredOneTimePasswords() {
	err := one_time_password.DeleteExpiredOneTimePasswords()
	if err != nil {
		fmt.Println(err)
	}
}
