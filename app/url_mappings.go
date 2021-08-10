package app

import (
	"github.com/StarsPoker/loginBackend/controllers/access_token"
	"github.com/StarsPoker/loginBackend/controllers/menus"
	"github.com/StarsPoker/loginBackend/controllers/profiles"
	"github.com/StarsPoker/loginBackend/controllers/users"
	"github.com/gin-gonic/gin"
)

func mapMenusRoutes(api *gin.RouterGroup) {
	menusGroup := api.Group("menus")

	menusGroup.GET("/menus/:profile_father", menus.MenusController.GetMenus)
	menusGroup.GET("/childrens", menus.MenusController.GetChildrens)
	menusGroup.POST("", menus.MenusController.InsertMenu)
	menusGroup.GET("/build_menu", menus.MenusController.BuildMenu)
	menusGroup.DELETE("/:menu_id", menus.MenusController.DeleteMenu)
	menusGroup.PUT("/update/:menu_id", menus.MenusController.UpdateMenu)
	menusGroup.PUT("/change_order_up/:menu_id", menus.MenusController.ChangeOrderUpMenu)
	menusGroup.PUT("/change_order_down/:menu_id", menus.MenusController.ChangeOrderDownMenu)
}

func mapProfileRoutes(api *gin.RouterGroup) {
	profileGroup := api.Group("profiles")

	profileGroup.GET("/:profile_id/profile_attendants/:search", profiles.ProfileController.GetProfileAttendants)
	profileGroup.GET("/:profile_id/usersadds", profiles.ProfileController.GetProfileUsersAdds)
	profileGroup.GET("/:profile_id/users", profiles.ProfileController.GetProfileUsers)
	profileGroup.GET("/:profile_id", profiles.ProfileController.GetProfile)
	profileGroup.GET("", profiles.ProfileController.GetProfiles)
	profileGroup.POST("", profiles.ProfileController.CreateProfile)
	profileGroup.POST("/user", profiles.ProfileController.CreateProfileUser)
	profileGroup.POST("/menu", profiles.ProfileController.CreateProfileMenu)
	profileGroup.POST("/menu_father", profiles.ProfileController.CreateProfileMenuFather)
	profileGroup.PUT("/updateProfile/:profile_id", profiles.ProfileController.UpdateProfile)
	profileGroup.PUT("/updateUser/:user_id", profiles.ProfileController.UpdateProfileUser)
	profileGroup.DELETE("/deleteUser/:user_id", profiles.ProfileController.DeleteProfileUser)
	profileGroup.DELETE("/deleteProfile/:profile_id", profiles.ProfileController.DeleteProfile)
	profileGroup.DELETE("/deleteMenu/:profilemenu_id", profiles.ProfileController.DeleteProfileMenu)
	profileGroup.DELETE("/deleteMenuFather/:profilemenu_id", profiles.ProfileController.DeleteProfileMenuFather)
}

func mapUsersRoutes(api *gin.RouterGroup) {
	usersGroup := api.Group("users")

	usersGroup.GET("", users.UserController.GetUsers)
	usersGroup.GET("/attendants/:search", users.UserController.GetAttendants)
	usersGroup.GET("/user/:user_id", users.UserController.GetUser)
	usersGroup.POST("", users.UserController.CreateUser)
	usersGroup.POST("/change_password/:user_id", users.UserController.ChangePassword)
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
	mapProfileRoutes(api)
	mapMenusRoutes(api)
}
