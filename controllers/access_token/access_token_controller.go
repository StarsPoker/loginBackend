package access_token

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/StarsPoker/loginBackend/domain/access_token"
	"github.com/StarsPoker/loginBackend/services"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
	"github.com/gin-gonic/gin"
)

var (
	AccessTokenController AccessTokenInterface = &accessTokenController{}
)

type AccessTokenInterface interface {
	GetById(c *gin.Context)
	Create(c *gin.Context)
	Delete(c *gin.Context)
	CheckAuth(c *gin.Context)
	DeleteExpiredTokens()
}

type accessTokenController struct {
}

func (cont *accessTokenController) GetById(c *gin.Context) {
	at, err := services.AccessTokenService.GetById(c.Param("access_token_id"))
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


	otp, err := services.AccessTokenService.Create(accessTokenRequest)

	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	fmt.Println("Created")

	c.JSON(http.StatusOK, otp)
}

func (cont *accessTokenController) Delete(c *gin.Context) {
	err := services.AccessTokenService.Delete(c.Param("access_token_id"))
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

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
