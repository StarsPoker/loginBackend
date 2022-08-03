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
	DeleteExpiredAccesTokens()
	DeleteExpiredOneTimePasswords()
}

func (s *accessTokenService) GetById(accessTokenId string) (*access_token.AccessToken, *rest_errors.RestErr) {

	accessTokenId = strings.TrimSpace(accessTokenId)
	if len(accessTokenId) == 0 {
		return nil, rest_errors.NewBadRequestError("invalid access token")
	}

	accessToken, err := access_token.GetById(accessTokenId)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func (s *accessTokenService) Create(accessTokenRequest access_token.AccessTokenRequest) (*one_time_password.OneTimePassword, *rest_errors.RestErr) {

	if err := accessTokenRequest.Validate(); err != nil {
		return nil, err
	}

	user := &users.User{
		Email:    accessTokenRequest.Username,
		Password: crypto_utils.GetMd5(accessTokenRequest.Password),
		Status:   users.StatusActive,
	}

	if err := user.FindByEmailAndPassword(); err != nil {
		return nil, err
	}

	var otp one_time_password.OneTimePassword
	otp = otp.CreateOtp(user)

	err := one_time_password.Insert(otp)

	if err != nil {
		return nil, err
	}

	if user.Contact != nil {
		go chat_repository.SendWhatsappMessage(otp, user)
	}
	go chat_repository.SendMail(otp, user)

	otp.Code = "anonimized"
	return &otp, nil
}

func (s *accessTokenService) ValidateAccessToken(accessTokenId string) *rest_errors.RestErr {
	if accessTokenId == "" {
		return rest_errors.NewUnauthorizedError("access token not found")
	}
	accessToken, err := AccessTokenService.GetById(accessTokenId)

	if err != nil {
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

		return rest_errors.NewUnauthorizedError("access token expired")
	}

	_ = access_token.UpdateLastInteraction(accessTokenId)

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

	if err := accessTokenRequest.Validate(); err != nil {
		return nil, err
	}

	user := &users.User{
		Email:    accessTokenRequest.Username,
		Password: crypto_utils.GetMd5(accessTokenRequest.Password),
		Status:   users.StatusActive,
	}

	if err := user.FindByEmailAndPassword(); err != nil {
		return nil, err
	}
	accessTokenRequest.ClientId = user.Id

	otp, err := one_time_password.GetAuth(accessTokenRequest)
	if err != nil {
		return nil, err
	}

	otp, err = one_time_password.CheckAuth(otp, accessTokenRequest)
	if err != nil {
		return nil, err
	}

	rec := false //recursive calling
	rec = one_time_password.DeleteAllById(*otp, rec)
	if !rec {
		return nil, nil
	}

	at := access_token.GetNewAccessToken(user.Id, user.Role)
	at.Generate()
	at.UserClientIp = client_ip
	err = access_token.Create(at)
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
