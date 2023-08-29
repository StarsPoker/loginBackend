package app

import (
	"github.com/StarsPoker/loginBackend/controllers/access_token"
	"github.com/StarsPoker/loginBackend/controllers/menus"
	"github.com/StarsPoker/loginBackend/controllers/profiles"
	"github.com/StarsPoker/loginBackend/controllers/routes"
	"github.com/StarsPoker/loginBackend/controllers/users"
	"github.com/gin-gonic/gin"
)

func mapMenusRoutes(api *gin.RouterGroup) {
	menusGroup := api.Group("menus")

	menusGroup.GET("/children_search/:search", menus.MenusController.GetChildrenSearch)
	menusGroup.GET("/menus/:profile_father", menus.MenusController.GetMenus)
	menusGroup.GET("/childrens", menus.MenusController.GetChildrens)
	menusGroup.GET("/profiles_relation/:menu_id", menus.MenusController.GetProfilesRelation)
	menusGroup.POST("", menus.MenusController.InsertMenu)
	menusGroup.GET("/build_menu", menus.MenusController.BuildMenu)
	menusGroup.GET("/build_menu_search/:menu_name", menus.MenusController.BuildMenuSearch)
	menusGroup.GET("/profile_permission/:menu_name", menus.MenusController.ProfilePermission)
	menusGroup.DELETE("/:menu_id", menus.MenusController.DeleteMenu)
	menusGroup.PUT("/update/:menu_id", menus.MenusController.UpdateMenu)
	menusGroup.PUT("/change_order_up/:menu_id", menus.MenusController.ChangeOrderUpMenu)
	menusGroup.PUT("/change_order_down/:menu_id", menus.MenusController.ChangeOrderDownMenu)
}

func mapRoutesRoutes(api *gin.RouterGroup) {
	routesGroup := api.Group("routes")

	routesGroup.GET("", routes.RouteController.GetRoutes)
	routesGroup.POST("", routes.RouteController.CreateRoute)
	routesGroup.PUT("/updateRoute/:route_id", routes.RouteController.UpdateRoute)
	routesGroup.DELETE("/deleteRoute/:route_id", routes.RouteController.DeleteRoute)
}

func mapProfileRoutes(api *gin.RouterGroup) {
	profileGroup := api.Group("profiles")

	profileGroup.GET("/:profile_id/profile_attendants/:search", profiles.ProfileController.GetProfileAttendants)
	profileGroup.GET("/:profile_id/profile_routesAdds/:search", profiles.ProfileController.GetProfileRoutesAdds)
	profileGroup.GET("/:profile_id/usersadds", profiles.ProfileController.GetProfileUsersAdds)
	profileGroup.GET("/:profile_id/permissions", profiles.ProfileController.GetProfilePermissions)
	profileGroup.GET("/:profile_id/users", profiles.ProfileController.GetProfileUsers)
	profileGroup.GET("/:profile_id/routes", profiles.ProfileController.GetProfileRoutes)
	profileGroup.GET("/:profile_id", profiles.ProfileController.GetProfile)
	profileGroup.GET("", profiles.ProfileController.GetProfiles)
	profileGroup.POST("", profiles.ProfileController.CreateProfile)
	profileGroup.POST("/user", profiles.ProfileController.CreateProfileUser)
	profileGroup.POST("/route", profiles.ProfileController.CreateProfileRoute)
	profileGroup.POST("/menu", profiles.ProfileController.CreateProfileMenu)
	profileGroup.POST("/menu_father", profiles.ProfileController.CreateProfileMenuFather)
	profileGroup.PUT("/updateProfile/:profile_id", profiles.ProfileController.UpdateProfile)
	profileGroup.PUT("/updateParam/:profile_id", profiles.ProfileController.UpdateParam)
	profileGroup.PUT("/updateUser/:user_id", profiles.ProfileController.UpdateProfileUser)
	profileGroup.DELETE("/deleteUser/:user_id", profiles.ProfileController.DeleteProfileUser)
	profileGroup.DELETE("/deleteRoute/:route_id", profiles.ProfileController.DeleteProfileRoute)
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
	usersGroup.PUT("/update/:user_id", users.UserController.UpdateUser)
	usersGroup.PUT("/update_useredit/:user_id", users.UserController.UpdateUserEdit)
	usersGroup.DELETE("/:user_id", users.UserController.DeleteUser)
}

func mapOauthRoutes(api *gin.RouterGroup) {
	oauthGroup := api.Group("oauth")
	// oauthGroup.GET("/access_token/:access_token_id", access_token.AccessTokenController.GetById)

	oauthGroup.POST("/qr_code_authenticator", access_token.AccessTokenController.GenerateQrCodeAuthenticator)
	oauthGroup.GET("/access_token", access_token.AccessTokenController.GetById)
	oauthGroup.POST("/access_token", access_token.AccessTokenController.Create)
	oauthGroup.DELETE("/access_token/:access_token_id", access_token.AccessTokenController.Delete)
	oauthGroup.POST("/access_token/check_auth", access_token.AccessTokenController.CheckAuth)
}

func mapUrls() {
	router.Use(CORSMiddleware())
	router.Use(OAuthMiddleware())

	api := router.Group("/api/auth/")

	api.OPTIONS("/*path", CORSMiddleware())

	mapUsersRoutes(api)
	mapOauthRoutes(api)
	mapProfileRoutes(api)
	mapMenusRoutes(api)
	mapRoutesRoutes(api)

	go access_token.AccessTokenController.DeleteExpiredTokens()
}
