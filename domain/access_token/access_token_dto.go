package access_token

import (
	"fmt"
	"strings"
	"time"

	"github.com/StarsPoker/loginBackend/utils/date_utils"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
	"gopkg.in/mgo.v2/bson"
)

const (
	expirationTime = 24
)

type AccessToken struct {
	Id          bson.ObjectId `json:"id" bson:"_id" `
	AccessToken string        `json:"access_token" bson:"access_token"`
	UserId      int64         `json:"user_id" bson:"user_id"`
	ClientId    int64         `json:"client_id" bson:"client_id"`
	Expires     int64         `json:"expires" bson:"expires"`
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

func GetNewAccessToken() AccessToken {
	return AccessToken{
		AccessToken: "brasil",
		UserId:      1,
		ClientId:    2,
		Expires:     date_utils.GetNow().Add(expirationTime * time.Hour).Unix(),
	}
}

func (at AccessToken) IsExpired() bool {
	now := date_utils.GetNow()
	expirationTime := time.Unix(at.Expires, 0)
	fmt.Println(expirationTime)
	return expirationTime.Before(now)
}
