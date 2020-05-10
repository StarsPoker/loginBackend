package services

import (
	"strings"

	"github.com/StarsPoker/loginBackend/domain/access_token"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
)

var (
	AccessTokenService AccessTokenServiceInterface = &accessTokenService{}
)

type accessTokenService struct {
}

type AccessTokenServiceInterface interface {
	GetById(string) (*access_token.AccessToken, *rest_errors.RestErr)
	Create() (*access_token.AccessToken, *rest_errors.RestErr)
	UpdateExpirationTime(*access_token.AccessToken) (*access_token.AccessToken, *rest_errors.RestErr)
}

func (s *accessTokenService) GetById(accessTokenId string) (*access_token.AccessToken, *rest_errors.RestErr) {

	accessTokenId = strings.TrimSpace(accessTokenId)
	if len(accessTokenId) == 0 {
		return nil, rest_errors.NewBadRequestError("invalid access token id")
	}

	accessToken, err := access_token.GetById(accessTokenId)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func (s *accessTokenService) Create() (*access_token.AccessToken, *rest_errors.RestErr) {
	at := access_token.GetNewAccessToken()

	err := access_token.Create(at)
	if err != nil {
		return nil, err
	}

	return &at, nil
}

func (s *accessTokenService) UpdateExpirationTime(at *access_token.AccessToken) (*access_token.AccessToken, *rest_errors.RestErr) {
	at, err := access_token.UpdateExpirationTime(at)
	if err != nil {
		return nil, err
	}

	return at, nil
}
