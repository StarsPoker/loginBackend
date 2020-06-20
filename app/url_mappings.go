package app

import (
	"github.com/StarsPoker/loginBackend/controllers/access_token"
	"github.com/StarsPoker/loginBackend/controllers/users"
	"github.com/gin-gonic/gin"
)

func mapUsersRoutes(api *gin.RouterGroup) {
	usersGroup := api.Group("users")

	usersGroup.GET("", users.UserController.GetUsers)
	usersGroup.GET("/attendances/:search", users.UserController.GetAttendances)
	usersGroup.GET("/user/:user_id", users.UserController.GetUser)
	usersGroup.POST("", users.UserController.CreateUser)
	usersGroup.PUT("/:user_id", users.UserController.UpdateUser)
	usersGroup.DELETE("/:user_id", users.UserController.DeleteUser)
}

func mapOauthRoutes(api *gin.RouterGroup) {
	oauthGroup := api.Group("oauth")
	oauthGroup.GET("/access_token/:access_token_id", access_token.AccessTokenController.GetById)
	oauthGroup.POST("/access_token", access_token.AccessTokenController.Create)
	oauthGroup.DELETE("/access_token/:access_token_id", access_token.AccessTokenController.Delete)
}

func mapUrls() {
	router.Use(CORSMiddleware())
	router.Use(OAuthMiddleware())

	api := router.Group("/api/")

	api.OPTIONS("/*path", CORSMiddleware())

	mapUsersRoutes(api)
	mapOauthRoutes(api)
}
