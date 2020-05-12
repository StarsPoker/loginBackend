package app

import (
	"github.com/StarsPoker/loginBackend/services"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
	"github.com/gin-gonic/gin"
)

const (
	loginRoute  = "oauth/accessToken"
	keyToken    = "oauth2_token"
	keyNextPage = "next"
)

func OAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		route := c.FullPath()
		if route != "/api/oauth/access_token" || c.Request.Method != "POST" {
			if len(c.Request.Header["Authorization"]) > 0 {
				token := c.Request.Header["Authorization"][0]
				if err := services.AccessTokenService.ValidateAccessToken(token); err != nil {
					c.JSON(err.Status, err)
					c.Abort()
					return
				}
				return
			}
			err := rest_errors.NewUnauthorizedError("access token not found")
			c.JSON(err.Status, err)
			c.Abort()
			return
		}
	}
}
