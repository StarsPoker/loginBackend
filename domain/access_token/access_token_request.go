package access_token

import (
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
)

const (
	grantTypePassword          = "password"
	grantTypeClientCredentials = "client_credentials"
)

type AccessTokenRequest struct {
	GrantType string `json:"grant_type"`
	Scope     string `json:"scope"`

	// Used for password grant type
	Username string `json:"username"`
	Password string `json:"password"`

	// Used for client_credentials grant type
	ClientId    int64  `json:"client_id"`
	ClientScret string `json:"client_secret"`

	// Used for grant otp searching
	ClientKey               string `json:"client_key"`
	UserIpFront             string `json:"user_ip_front"`
	AuthenticatorConfigured bool   `json:"authenticator_configured"`
}

func (at *AccessTokenRequest) Validate() *rest_errors.RestErr {

	if at.GrantType != grantTypePassword && at.GrantType != grantTypeClientCredentials {
		return rest_errors.NewBadRequestError("invalid grant_type parameter")
	}

	switch at.GrantType {
	case grantTypePassword:
		if at.Username == "" {
			return rest_errors.NewBadRequestError("Usuário deve ser informado")
		}

		if at.Password == "" {
			return rest_errors.NewBadRequestError("Senha deve ser informada")
		}
		break
	case grantTypeClientCredentials:
		if at.ClientId == 0 {
			return rest_errors.NewBadRequestError("Id do cliente deve ser informado")
		}

		if at.ClientScret == "" {
			return rest_errors.NewBadRequestError("Senha deve ser informada")
		}
	}

	return nil
}
