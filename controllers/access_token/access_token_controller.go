package access_token

import (
	"net/http"

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

	at, err := services.AccessTokenService.Create(accessTokenRequest)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, at)
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
