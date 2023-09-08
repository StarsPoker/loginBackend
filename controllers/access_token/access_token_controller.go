package access_token

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/StarsPoker/loginBackend/domain/access_token"
	"github.com/StarsPoker/loginBackend/services"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
	"github.com/gin-gonic/gin"
)

var (
	stars_env                                  = os.Getenv("stars_env")
	AccessTokenController AccessTokenInterface = &accessTokenController{}
)

type AccessTokenInterface interface {
	GetById(c *gin.Context)
	Create(c *gin.Context)
	Delete(c *gin.Context)
	CheckAuth(c *gin.Context)
	DeleteExpiredTokens()
	GenerateQrCodeAuthenticator(c *gin.Context)
	ValidateAccessToken(c *gin.Context)
}

type accessTokenController struct {
}

func (cont *accessTokenController) GetById(c *gin.Context) {
	access_token_id := c.Request.Header["Authorization"][0]
	if strings.Contains(access_token_id, "Bearer") {
		access_token_id = access_token_id[7:]
	}
	at, err := services.AccessTokenService.GetById(access_token_id)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, at)
}

func (cont *accessTokenController) Create(c *gin.Context) {
	var accessTokenRequest access_token.AccessTokenRequest

	if err := c.ShouldBindJSON(&accessTokenRequest); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	if stars_env == "development" {
		host := c.Request.Host
		client_ip := c.ClientIP()
		otp, err := services.AccessTokenService.CreateDevelopment(accessTokenRequest, host, client_ip)
		if err != nil {
			c.JSON(err.Status, err)
			return
		}
		http.SetCookie(c.Writer, &http.Cookie{
			Name:    "sx_access_token",
			Value:   otp.AccessToken.AccessToken,
			Expires: otp.AccessToken.ExpirationTime,
		})
		c.JSON(http.StatusOK, otp)
	} else {
		otp, err := services.AccessTokenService.Create(accessTokenRequest)
		if err != nil {
			fmt.Println("Created")
			c.JSON(err.Status, err)
			return
		}
		http.SetCookie(c.Writer, &http.Cookie{
			Name:    "sx_access_token",
			Value:   otp.AccessToken.AccessToken,
			Expires: otp.AccessToken.ExpirationTime,
		})
		c.JSON(http.StatusOK, otp)
	}
}

func (cont *accessTokenController) Delete(c *gin.Context) {
	err := services.AccessTokenService.Delete(c.Param("access_token_id"))
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:  "sx_access_token",
		Value: "",
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (cont *accessTokenController) CheckAuth(c *gin.Context) {
	var accessTokenRequest access_token.AccessTokenRequest

	if err := c.ShouldBindJSON(&accessTokenRequest); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	host := c.Request.Host
	client_ip := c.ClientIP()

	otp, err := services.AccessTokenService.CheckAuth(accessTokenRequest, host, client_ip)
	if err != nil {
		c.JSON(err.Status, err)
		return

	}

	c.JSON(http.StatusOK, otp)
}

func (cont *accessTokenController) DeleteExpiredTokens() {

	services.AccessTokenService.DeleteExpiredAccesTokens()
	services.AccessTokenService.DeleteExpiredOneTimePasswords()

}

func (cont *accessTokenController) GenerateQrCodeAuthenticator(c *gin.Context) {
	var accessTokenRequest access_token.AccessTokenRequest

	if err := c.ShouldBindJSON(&accessTokenRequest); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	qrCode, err := services.AccessTokenService.GenerateQrCodeAuthenticator(accessTokenRequest)
	if err != nil {
		c.JSON(err.Status, err)
		return

	}

	c.JSON(http.StatusOK, qrCode)
}

func (cont *accessTokenController) ValidateAccessToken(c *gin.Context) {
	token := ""
	cookie, errGetCookie := c.Cookie("sx_access_token")
	setCookie := false
	if errGetCookie != nil || cookie == "" {
		if len(c.Request.Header["Authorization"]) == 0 {
			c.JSON(http.StatusUnauthorized, map[string]string{"success": "unauthorized"})
			return
		}
		token = c.Request.Header["Authorization"][0]
	} else {
		token = cookie
	}

	at, err := services.AccessTokenService.GetById(token)

	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	if setCookie {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:    "sx_access_token",
			Value:   at.AccessToken,
			Expires: at.ExpirationTime,
		})
	}

	c.JSON(http.StatusOK, at)
}
