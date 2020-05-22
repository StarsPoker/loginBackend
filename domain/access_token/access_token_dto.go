package access_token

import (
	"fmt"
	"strings"
	"time"

	"github.com/StarsPoker/loginBackend/utils/crypto_utils.go"
	"github.com/StarsPoker/loginBackend/utils/date_utils"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
)

const (
	expirationTime = 12
)

type AccessToken struct {
	AccessToken     string `json:"access_token" bson:"access_token"`
	Role            int64  `json:"role" bson:"role"`
	UserId          int64  `json:"user_id" bson:"user_id"`
	ClientId        int64  `json:"client_id" bson:"client_id"`
	Expires         int64  `json:"expires" bson:"expires"`
	LastInteraction int64  `json:"last_interaction" bson:"last_interaction"`
	Status          int64  `json:"status" bson:"status"`
}

func (at *AccessToken) Validate() *rest_errors.RestErr {
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if at.AccessToken == "" {
		return rest_errors.NewBadRequestError("invalid access token id")
	}

	if at.UserId <= 0 {
		return rest_errors.NewBadRequestError("invalid user id")
	}

	if at.ClientId <= 0 {
		return rest_errors.NewBadRequestError("invalid client id")
	}

	if at.Expires <= 0 {
		return rest_errors.NewBadRequestError("invalid expiration time")
	}
	return nil
}

func GetNewAccessToken(userId int64, role int64) AccessToken {
	return AccessToken{
		Role:            role,
		UserId:          userId,
		Expires:         date_utils.GetNow().Add(expirationTime * time.Hour).Unix(),
		LastInteraction: date_utils.GetNow().Unix(),
		Status:          1,
	}
}

func (at AccessToken) IsExpired() bool {
	now := date_utils.GetNow()
	expirationTime := time.Unix(at.Expires, 0)

	return expirationTime.Before(now)
}

func (at *AccessToken) Generate() {
	at.AccessToken = crypto_utils.GetMd5(fmt.Sprintf("at-%d-%d-ran", at.UserId, at.Expires))
}
