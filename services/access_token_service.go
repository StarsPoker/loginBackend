package services

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/StarsPoker/loginBackend/domain/access_token"
	"github.com/StarsPoker/loginBackend/domain/one_time_password"
	"github.com/StarsPoker/loginBackend/domain/users"
	"github.com/StarsPoker/loginBackend/utils/crypto_utils"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
	"github.com/skip2/go-qrcode"
	"github.com/xlzd/gotp"
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
	GenerateQrCodeAuthenticator(accessTokenRequest access_token.AccessTokenRequest) (*access_token.QrCodeAuthenticator, *rest_errors.RestErr)
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
	otp = otp.CreateOtp(user, false)
	if !user.AuthenticatorConfigured {
		qr, errQr := GenerateQrCodeAuthenticatorByOtp(otp)
		if errQr != nil {
			return nil, errQr
		}
		otp.QrCode = *qr
	}
	otp.AuthenticatorConfigured = user.AuthenticatorConfigured
	err := one_time_password.Insert(otp)

	if err != nil {
		return nil, err
	}

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
		return rest_errors.NewUnauthorizedError("access token expired")
	}

	_ = access_token.UpdateLastInteraction(accessTokenId)
	return nil
}

func (s *accessTokenService) Delete(accessTokenId string) *rest_errors.RestErr {
	access_token.Delete(accessTokenId)

	return nil
}

func (s *accessTokenService) CheckAuth(accessTokenRequest access_token.AccessTokenRequest, host string, client_ip string) (*one_time_password.OneTimePassword, *rest_errors.RestErr) {

	otp, err := one_time_password.GetAuth(accessTokenRequest)
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
	if user.OTPSecret == nil {
		return nil, rest_errors.NewInternalServerError("Reinicie o fluxo de login")
	}
	totp := gotp.NewDefaultTOTP(*user.OTPSecret)
	ok := totp.Verify(accessTokenRequest.ClientScret, time.Now().Unix())
	if ok {
		rec := false //recursive calling
		rec = one_time_password.DeleteAllById(*otp, rec)
		if !rec {
			return nil, nil
		}

		at := access_token.GetNewAccessToken(user.Id, *user.Role, user.Inscription)
		at.Generate()
		at.UserHost = host
		at.UserClientIp = client_ip
		at.UserIpFront = accessTokenRequest.UserIpFront
		err = access_token.Create(at)
		if err != nil {
			return nil, err
		}
		otp.AccessToken = at
		if !user.AuthenticatorConfigured {
			user.AuthenticatorConfigured = true
			user.Update()
		}
	} else {
		return nil, rest_errors.NewInternalServerError("Código inválido ou expirado")
	}

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
	at := access_token.GetNewAccessToken(user.Id, *user.Role, user.Inscription)
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

func GenerateQrCodeAuthenticatorByOtp(otp one_time_password.OneTimePassword) (*access_token.QrCodeAuthenticator, *rest_errors.RestErr) {
	user, errGetUser := UsersService.GetUser(otp.UserId)
	if errGetUser != nil {
		return nil, errGetUser
	}
	secret := ""
	if user.OTPSecret == nil {
		newSecret := gotp.RandomSecret(32)
		user.OTPSecret = &newSecret

		errUpdateUser := user.Update()
		if errUpdateUser != nil {
			return nil, errUpdateUser
		}
		secret = newSecret
	} else {
		secret = *user.OTPSecret
	}
	totp := gotp.NewDefaultTOTP(secret)
	uri := totp.ProvisioningUri(user.Email, "GrupoSX")

	qrCodeImage, errGenQrImage := generateQR(uri)
	if errGenQrImage != nil {
		return nil, errGenQrImage
	}
	qr := &access_token.QrCodeAuthenticator{
		URI:    uri,
		Base64: qrCodeImage,
	}
	return qr, nil
}

func (s *accessTokenService) GenerateQrCodeAuthenticator(accessTokenRequest access_token.AccessTokenRequest) (*access_token.QrCodeAuthenticator, *rest_errors.RestErr) {
	otp, err := one_time_password.GetAuth(accessTokenRequest)
	if err != nil {
		return nil, err
	}

	return GenerateQrCodeAuthenticatorByOtp(*otp)
}

func generateQR(url string) (string, *rest_errors.RestErr) {
	qrCode, _ := qrcode.New(url, qrcode.Medium)
	png, err := qrCode.PNG(256)
	if err != nil {
		return "", rest_errors.NewInternalServerError("error in generating qr code")
	}
	sEnc := base64.StdEncoding.EncodeToString(png)
	return sEnc, nil
}
