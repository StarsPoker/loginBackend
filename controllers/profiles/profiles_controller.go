package profiles

import (
	"net/http"
	"strconv"

	"github.com/StarsPoker/loginBackend/domain/profiles"
	"github.com/StarsPoker/loginBackend/services"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
	"github.com/gin-gonic/gin"
)

var (
	ProfileController ProfileInterface = &profileController{}
)

type ProfileInterface interface {
	CreateProfile(c *gin.Context)
	CreateProfileUser(c *gin.Context)
	CreateProfileRoute(c *gin.Context)
	CreateProfileMenu(c *gin.Context)
	CreateProfileMenuFather(c *gin.Context)
	GetProfiles(c *gin.Context)
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
	UpdateParam(c *gin.Context)
	DeleteProfile(c *gin.Context)
	DeleteProfileUser(c *gin.Context)
	DeleteProfileRoute(c *gin.Context)
	DeleteProfileMenu(c *gin.Context)
	DeleteProfileMenuFather(c *gin.Context)
	GetProfileUsers(c *gin.Context)
	GetProfileRoutes(c *gin.Context)
	GetProfileUsersAdds(c *gin.Context)
	GetProfilePermissions(c *gin.Context)
	GetProfileAttendants(c *gin.Context)
	GetProfileRoutesAdds(c *gin.Context)
	UpdateProfileUser(c *gin.Context)
	GetProfileUser(c *gin.Context)
	GetProfileValidation(c *gin.Context)
}

type profileController struct {
}

func (cont *profileController) getProfileId(profileIdParam string) (int64, *rest_errors.RestErr) {
	profileId, profileErr := strconv.ParseInt(profileIdParam, 10, 64)
	if profileErr != nil {
		return 0, rest_errors.NewBadRequestError("profile id should be a number")
	}
	return profileId, nil
}

func (cont *profileController) getProfileUserId(profileUserIdParam string) (int64, *rest_errors.RestErr) {
	profileUserId, profileErr := strconv.ParseInt(profileUserIdParam, 10, 64)
	if profileErr != nil {
		return 0, rest_errors.NewBadRequestError("profile_user id should be a number")
	}
	return profileUserId, nil
}

func (cont *profileController) GetProfile(c *gin.Context) {

	profileId, idErr := cont.getProfileId(c.Param("profile_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	withdrawal, getErr := services.ProfilesService.GetProfile(profileId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}

	c.JSON(http.StatusOK, withdrawal)
}

func (cont *profileController) GetProfileValidation(c *gin.Context) {

	// profileId, idErr := cont.getProfileId(c.Param("profile_id"))
	// if idErr != nil {
	// 	c.JSON(idErr.Status, idErr)
	// 	return
	// }

	validation, getErr := services.ProfilesService.GetProfileValidation(c.Param("profile_id"), c.Param("param_to_check"))
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}

	c.JSON(http.StatusOK, validation)
}

func (cont *profileController) GetProfileUser(c *gin.Context) {

	profileUserId, idErr := cont.getProfileUserId(c.Param("profile_user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	withdrawal, getErr := services.ProfilesService.GetProfileUser(profileUserId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}

	c.JSON(http.StatusOK, withdrawal)
}

func buildFilter(c *gin.Context) *profiles.Filter {
	var filter profiles.Filter
	filter.Id = c.Query("id")
	filter.Name = c.Query("name")
	filter.ProfileCode = c.Query("profile_code")
	return &filter
}

func (cont *profileController) GetProfiles(c *gin.Context) {

	pageParam := c.Query("page")
	itemsPerPageParam := c.Query("itemsPerPage")
	filter := buildFilter(c)
	page := 1
	itemsPerPage := 10

	token := c.Request.Header["Authorization"][0]
	at, tokenErr := services.AccessTokenService.GetById(token)
	if tokenErr != nil {
		err := rest_errors.NewBadRequestError("Invalid access token")
		c.JSON(err.Status, err)
		return
	}
	user, _ := services.UsersService.GetUser(at.UserId)

	var userId = user.Id

	if pageParam != "" {
		page, _ = strconv.Atoi(pageParam)
	}

	if itemsPerPageParam != "" {
		itemsPerPage, _ = strconv.Atoi(itemsPerPageParam)
	}

	result, total, getErr := services.ProfilesService.GetProfiles(page, itemsPerPage, filter, userId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}

	var profilesResponse profiles.ProfilesResponse

	profilesResponse.Total = *total
	profilesResponse.Profiles = result

	c.JSON(http.StatusOK, profilesResponse)
}

func (cont *profileController) CreateProfile(c *gin.Context) {
	var w profiles.Profile
	if err := c.ShouldBindJSON(&w); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	result, saveErr := services.ProfilesService.CreateProfile(w)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (cont *profileController) UpdateProfile(c *gin.Context) {
	profileIdParam := c.Param("profile_id")
	profileId, profileErr := strconv.ParseInt(profileIdParam, 10, 64)

	if profileErr != nil {
		restErr := rest_errors.NewBadRequestError("profile id should be a number")
		c.JSON(restErr.Status, restErr)
		return
	}

	var b profiles.Profile
	if err := c.ShouldBindJSON(&b); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	b.Id = profileId

	result, updateErr := services.ProfilesService.UpdateProfile(b)
	if updateErr != nil {
		c.JSON(updateErr.Status, updateErr)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (cont *profileController) UpdateParam(c *gin.Context) {
	profileIdParam := c.Param("profile_id")
	profileId, profileErr := strconv.ParseInt(profileIdParam, 10, 64)

	if profileErr != nil {
		restErr := rest_errors.NewBadRequestError("profile id should be a number")
		c.JSON(restErr.Status, restErr)
		return
	}

	var b profiles.Profile
	if err := c.ShouldBindJSON(&b); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	b.Id = profileId

	result, updateErr := services.ProfilesService.UpdateParam(b)
	if updateErr != nil {
		c.JSON(updateErr.Status, updateErr)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (cont *profileController) DeleteProfile(c *gin.Context) {
	profileIdParam := c.Param("profile_id")

	profileId, profileErr := strconv.ParseInt(profileIdParam, 10, 64)

	if profileErr != nil {
		restErr := rest_errors.NewBadRequestError("profile id should be a number")
		c.JSON(restErr.Status, restErr)
		return
	}

	var profileDelete profiles.Profile

	profileDelete.Id = profileId

	deleteErr := services.ProfilesService.DeleteProfile(profileDelete)
	if deleteErr != nil {
		c.JSON(deleteErr.Status, deleteErr)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func (cont *profileController) DeleteProfileUser(c *gin.Context) {
	userIdParam := c.Param("user_id")

	userId, userErr := strconv.ParseInt(userIdParam, 10, 64)

	if userErr != nil {
		restErr := rest_errors.NewBadRequestError("user id should be a number")
		c.JSON(restErr.Status, restErr)
		return
	}

	var userDelete profiles.ProfileUser
	userDelete.Id = userId

	deleteErr := services.ProfilesService.DeleteProfileUser(userDelete)
	if deleteErr != nil {
		c.JSON(deleteErr.Status, deleteErr)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func (cont *profileController) DeleteProfileRoute(c *gin.Context) {
	routeIdParam := c.Param("route_id")

	routeId, routeErr := strconv.ParseInt(routeIdParam, 10, 64)

	if routeErr != nil {
		restErr := rest_errors.NewBadRequestError("route id should be a number")
		c.JSON(restErr.Status, restErr)
		return
	}

	var routeDelete profiles.ProfileRoute
	routeDelete.Id = routeId

	deleteErr := services.ProfilesService.DeleteProfileRoute(routeDelete)
	if deleteErr != nil {
		c.JSON(deleteErr.Status, deleteErr)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func (cont *profileController) DeleteProfileMenu(c *gin.Context) {
	menuIdParam := c.Param("profilemenu_id")
	menuId, menuErr := strconv.ParseInt(menuIdParam, 10, 64)

	if menuErr != nil {
		restErr := rest_errors.NewBadRequestError("menu id should be a number")
		c.JSON(restErr.Status, restErr)
		return
	}

	var menuDelete profiles.ProfileMenu
	menuDelete.Id = menuId

	deleteErr := services.ProfilesService.DeleteProfileMenu(menuDelete)
	if deleteErr != nil {
		c.JSON(deleteErr.Status, deleteErr)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func (cont *profileController) DeleteProfileMenuFather(c *gin.Context) {
	menuIdParam := c.Param("profilemenu_id")
	menuId, menuErr := strconv.ParseInt(menuIdParam, 10, 64)

	if menuErr != nil {
		restErr := rest_errors.NewBadRequestError("menu id should be a number")
		c.JSON(restErr.Status, restErr)
		return
	}

	var menuDelete profiles.ProfileMenu
	menuDelete.Id = menuId

	deleteErr := services.ProfilesService.DeleteProfileMenuFather(menuDelete)
	if deleteErr != nil {
		c.JSON(deleteErr.Status, deleteErr)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func (cont *profileController) GetProfileUsers(c *gin.Context) {
	profileIdParam := c.Param("profile_id")

	profileId, profileErr := strconv.ParseInt(profileIdParam, 10, 64)

	if profileErr != nil {
		restErr := rest_errors.NewBadRequestError("profile id should be a number")
		c.JSON(restErr.Status, restErr)
		return
	}

	pageParam := c.Query("page")
	itemsPerPageParam := c.Query("itemsPerPage")
	filter := buildFilter(c)
	page := 1
	itemsPerPage := 10

	if pageParam != "" {
		page, _ = strconv.Atoi(pageParam)
	}

	if itemsPerPageParam != "" {
		itemsPerPage, _ = strconv.Atoi(itemsPerPageParam)
	}

	result, total, getErr := services.ProfilesService.GetProfileUsers(page, itemsPerPage, filter, profileId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}

	var usersResponse profiles.UsersResponse

	usersResponse.Total = *total
	usersResponse.Users = result

	c.JSON(http.StatusOK, usersResponse)
}

func (cont *profileController) GetProfileRoutes(c *gin.Context) {
	profileIdParam := c.Param("profile_id")

	profileId, profileErr := strconv.ParseInt(profileIdParam, 10, 64)

	if profileErr != nil {
		restErr := rest_errors.NewBadRequestError("profile id should be a number")
		c.JSON(restErr.Status, restErr)
		return
	}

	pageParam := c.Query("page")
	itemsPerPageParam := c.Query("itemsPerPage")
	filter := buildFilter(c)
	page := 1
	itemsPerPage := 10

	if pageParam != "" {
		page, _ = strconv.Atoi(pageParam)
	}

	if itemsPerPageParam != "" {
		itemsPerPage, _ = strconv.Atoi(itemsPerPageParam)
	}

	result, total, getErr := services.ProfilesService.GetProfileRoutes(page, itemsPerPage, filter, profileId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}

	var routesResponse profiles.RoutesResponse

	routesResponse.Total = *total
	routesResponse.Routes = result

	c.JSON(http.StatusOK, routesResponse)
}

func (cont *profileController) GetProfileUsersAdds(c *gin.Context) {
	profileIdParam := c.Param("profile_id")

	profileId, profileErr := strconv.ParseInt(profileIdParam, 10, 64)

	if profileErr != nil {
		restErr := rest_errors.NewBadRequestError("profile id should be a number")
		c.JSON(restErr.Status, restErr)
		return
	}

	pageParam := c.Query("page")
	itemsPerPageParam := c.Query("itemsPerPage")
	filter := buildFilter(c)
	page := 1
	itemsPerPage := 10

	if pageParam != "" {
		page, _ = strconv.Atoi(pageParam)
	}

	if itemsPerPageParam != "" {
		itemsPerPage, _ = strconv.Atoi(itemsPerPageParam)
	}

	result, getErr := services.ProfilesService.GetProfileUsersAdds(page, itemsPerPage, filter, profileId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}

	var usersResponse profiles.UsersResponse

	usersResponse.Users = result

	c.JSON(http.StatusOK, usersResponse)
}

func (cont *profileController) GetProfilePermissions(c *gin.Context) {

	result, getErr := services.ProfilesService.GetProfilePermissions(c.Param("profile_id"))
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (cont *profileController) GetProfileAttendants(c *gin.Context) {
	profileIdParam := c.Param("profile_id")

	profileId, profileErr := strconv.ParseInt(profileIdParam, 10, 64)

	if profileErr != nil {
		restErr := rest_errors.NewBadRequestError("profile id should be a number")
		c.JSON(restErr.Status, restErr)
		return
	}

	search := c.Param("search")

	userList, getErr := services.ProfilesService.GetProfileAttendants(search, profileId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	c.JSON(http.StatusOK, userList)
}

func (cont *profileController) GetProfileRoutesAdds(c *gin.Context) {
	profileIdParam := c.Param("profile_id")

	profileId, profileErr := strconv.ParseInt(profileIdParam, 10, 64)

	if profileErr != nil {
		restErr := rest_errors.NewBadRequestError("profile id should be a number")
		c.JSON(restErr.Status, restErr)
		return
	}

	search := c.Param("search")

	routeList, getErr := services.ProfilesService.GetProfileRoutesAdds(search, profileId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	c.JSON(http.StatusOK, routeList)
}

func (cont *profileController) CreateProfileUser(c *gin.Context) {
	var p profiles.ProfileUser

	if err := c.ShouldBindJSON(&p); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	result, saveErr := services.ProfilesService.CreateProfileUser(p)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (cont *profileController) CreateProfileRoute(c *gin.Context) {
	var p profiles.ProfileRoute

	if err := c.ShouldBindJSON(&p); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	result, saveErr := services.ProfilesService.CreateProfileRoute(p)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (cont *profileController) CreateProfileMenu(c *gin.Context) {
	var p profiles.ProfileMenu

	if err := c.ShouldBindJSON(&p); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	result, saveErr := services.ProfilesService.CreateProfileMenu(p)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (cont *profileController) CreateProfileMenuFather(c *gin.Context) {
	var p profiles.ProfileMenu

	if err := c.ShouldBindJSON(&p); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	result, saveErr := services.ProfilesService.CreateProfileMenuFather(p)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (cont *profileController) UpdateProfileUser(c *gin.Context) {
	userIdParam := c.Param("user_id")
	userId, userErr := strconv.ParseInt(userIdParam, 10, 64)

	if userErr != nil {
		restErr := rest_errors.NewBadRequestError("user id should be a number")
		c.JSON(restErr.Status, restErr)
		return
	}

	var userUpdate profiles.ProfileUser
	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	userUpdate.Id = userId

	result, updateErr := services.ProfilesService.UpdateProfileUser(userUpdate)
	if updateErr != nil {
		c.JSON(updateErr.Status, updateErr)
		return
	}
	c.JSON(http.StatusOK, result)
}
