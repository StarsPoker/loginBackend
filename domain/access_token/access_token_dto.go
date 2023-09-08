package access_token

import (
	"fmt"
	"os"
	"strings"
	"time"

	// "github.com/StarsPoker/loginBackend/utils/crypto_utils"
	"github.com/StarsPoker/loginBackend/utils/date_utils"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
	"github.com/golang-jwt/jwt/v4"
)

const (
	expirationTime = 12
	jwt_key_env    = "jwt_key"
)

var jwtKey = []byte(os.Getenv(jwt_key_env))

type AccessToken struct {
	AccessToken     string    `json:"access_token" bson:"access_token"`
	Role            int64     `json:"role" bson:"role"`
	UserId          int64     `json:"user_id" bson:"user_id"`
	ClientId        int64     `json:"client_id" bson:"client_id"`
	Expires         int64     `json:"expires" bson:"expires"`
	LastInteraction int64     `json:"last_interaction" bson:"last_interaction"`
	Status          int64     `json:"status" bson:"status"`
	UserHost        string    `json:"user_host" bson:"user_host"`
	UserClientIp    string    `json:"user_client_ip" bson:"user_client_ip"`
	UserIpFront     string    `json:"user_ip_front" bson:"user_ip_front"`
	Doc             *string   `json:"doc" bson:"doc"`
	ExpirationTime  time.Time `json:"expiration_time" bson:"expiration_time"`
	jwt.RegisteredClaims
}

type QrCodeAuthenticator struct {
	URI        string `json:"uri"`
	Base64     string `json:"base64"`
	ShowQrCode bool   `json:"show_qr_code"`
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

func GetNewAccessToken(userId int64, role int64, doc *string) AccessToken {
	return AccessToken{
		Role:            role,
		UserId:          userId,
		Expires:         date_utils.GetNow().Add(expirationTime * time.Hour).Unix(),
		LastInteraction: date_utils.GetNow().Unix(),
		Status:          1,
		Doc:             doc,
	}
}

func (at AccessToken) IsExpired() bool {
	if time.Until(at.RegisteredClaims.ExpiresAt.Time) < 0 {
		return true
	}
	return false

}

func (at *AccessToken) Generate() {
	time.Local, _ = time.LoadLocation("America/Sao_Paulo")
	expirationTimeJwt := time.Now().Add(expirationTime * time.Hour)

	at.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expirationTimeJwt),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, at)
	fmt.Println(string(jwtKey))
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		fmt.Println(err)
		return
	}
	at.Expires = expirationTimeJwt.Unix()
	at.ExpirationTime = expirationTimeJwt
	at.AccessToken = tokenString
}

func CheckToken(accessTokenId string) (*AccessToken, *rest_errors.RestErr) {
	claims := &AccessToken{}

	fmt.Println(string(jwtKey))
	tkn, err := jwt.ParseWithClaims(accessTokenId, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		fmt.Println(err)
		return nil, rest_errors.NewUnauthorizedError("invalid access token")
	}

	if !tkn.Valid {
		fmt.Println(err)
		return nil, rest_errors.NewUnauthorizedError("invalid access token")
	}
	at := &AccessToken{
		Role:            claims.Role,
		UserId:          claims.UserId,
		LastInteraction: date_utils.GetNow().Unix(),
		Status:          1,
		ClientId:        claims.ClientId,
		AccessToken:     accessTokenId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: claims.RegisteredClaims.ExpiresAt,
		},
	}
	return at, nil
}
