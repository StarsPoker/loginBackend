package one_time_password

import (
	"fmt"

	"github.com/StarsPoker/loginBackend/domain/access_token"
	"github.com/StarsPoker/loginBackend/domain/users"
	"github.com/StarsPoker/loginBackend/utils/crypto_utils"
	"github.com/StarsPoker/loginBackend/utils/date_utils"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
)

type OneTimePassword struct {
	Key         string                   `json:"key" bson:"key"`
	Code        string                   `json:"code" bson:"code"`
	UserId      int64                    `json:"user_id" bson:"user_id"`
	Expires     int64                    `json:"expires" bson:"expires"`
	AccessToken access_token.AccessToken `json:"access_token" bson:"access_token"`
}

const (
	MIN_TOKEN      = 1000
	MAX_TOKEN      = 9999
	TIME_TO_EXPIRE = 300
)

func (otp *OneTimePassword) CreateOtp(user *users.User) OneTimePassword {
	otp.Expires = date_utils.GetNow().Unix() + TIME_TO_EXPIRE
	otp.Code = crypto_utils.GetToken(MIN_TOKEN, MAX_TOKEN)
	otp.Key = crypto_utils.GetMd5(fmt.Sprintf("at-%d-%s-ran", user.Id, otp.Code))
	otp.UserId = user.Id
	return *otp
}

func CheckAuth(otp *OneTimePassword, at access_token.AccessTokenRequest) (*OneTimePassword, *rest_errors.RestErr) {
	if otp.Code != at.ClientScret {
		return nil, rest_errors.NewInternalServerError("invalid Token provided")
	}
	if otp.Expires > date_utils.GetNow().Unix() {
		return nil, rest_errors.NewInternalServerError("Token expired")
	}
	return otp, nil
}
