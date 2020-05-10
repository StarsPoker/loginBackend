package access_token

import (
	"net/http"

	"github.com/StarsPoker/loginBackend/domain/access_token"
	"github.com/StarsPoker/loginBackend/services"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
	"github.com/gin-gonic/gin"
)

func GetById(c *gin.Context) {
	at, err := services.AccessTokenService.GetById(c.Param("access_token_id"))
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, at)
}

func Create(c *gin.Context) {
	at, err := services.AccessTokenService.Create()
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, at)
}

func UpdateExpirationTime(c *gin.Context) {

	var at access_token.AccessToken
	if err := c.ShouldBindJSON(&at); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	result, saveErr := services.AccessTokenService.UpdateExpirationTime(&at)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result)
}
