package app

import (
	"strings"

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
		route := c.Request.URL.Path
		if !strings.Contains(route, "/oauth/access_token") && !strings.Contains(route, "qr_code_authenticator") {
			cookie, errGetCookie := c.Cookie("sx_access_token")

			if errGetCookie == nil && cookie != "" {
				token := cookie
				if err := services.AccessTokenService.ValidateAccessToken(token); err != nil {
					c.JSON(err.Status, err)
					c.Abort()
					return
				}
				return
			} else if len(c.Request.Header["Authorization"]) > 0 {
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
