package users

import (
	"net/http"
	"strconv"

	"github.com/StarsPoker/loginBackend/domain/users"
	"github.com/StarsPoker/loginBackend/services"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
	"github.com/gin-gonic/gin"
)

var (
	UserController UserControllerInterface = &userController{}
)

type UserControllerInterface interface {
	getUserId(string) (int64, *rest_errors.RestErr)
	CreateUser(*gin.Context)
	UpdateUser(*gin.Context)
	DeleteUser(c *gin.Context)
	GetUser(c *gin.Context)
	GetUsers(c *gin.Context)
	GetAttendants(c *gin.Context)
	ChangePassword(c *gin.Context)
}

type userController struct {
}

func (cont *userController) getUserId(userIdParam string) (int64, *rest_errors.RestErr) {
	userId, userErr := strconv.ParseInt(userIdParam, 10, 64)
	if userErr != nil {
		return 0, rest_errors.NewBadRequestError("user id should be a number")
	}
	return userId, nil
}

func (cont *userController) CreateUser(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	result, saveErr := services.UsersService.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func (cont *userController) UpdateUser(c *gin.Context) {
	userId, idErr := UserController.getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	user.Id = userId

	result, updateErr := services.UsersService.UpdateUser(user)
	if updateErr != nil {
		c.JSON(updateErr.Status, updateErr)
		return
	}

	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func (cont *userController) ChangePassword(c *gin.Context) {
	userId, idErr := UserController.getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	var user users.ChangePassword
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	user.Id = userId

	errChange := services.UsersService.ChangePassword(user)

	if errChange != nil {
		c.JSON(errChange.Status, errChange)
		return
	}

	c.JSON(http.StatusOK, map[string]string{"status": "changed"})
}

func (cont *userController) DeleteUser(c *gin.Context) {
	userId, idErr := UserController.getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	var user users.User

	user.Id = userId

	deleteErr := services.UsersService.DeleteUser(user)
	if deleteErr != nil {
		c.JSON(deleteErr.Status, deleteErr)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func (cont *userController) GetUser(c *gin.Context) {
	userId, idErr := UserController.getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	user, getErr := services.UsersService.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}

	c.JSON(http.StatusOK, user.Marshall(c.GetHeader("X-Public") == "true"))
}

func (cont *userController) GetUsers(c *gin.Context) {

	pageParam := c.Query("page")
	itemsPerPageParam := c.Query("itemsPerPage")
	page := 1
	itemsPerPage := 10

	if pageParam != "" {
		page, _ = strconv.Atoi(pageParam)
	}

	if itemsPerPageParam != "" {
		itemsPerPage, _ = strconv.Atoi(itemsPerPageParam)
	}

	result, total, getErr := services.UsersService.GetUsers(page, itemsPerPage)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}

	users_result := result.Marshall(c.GetHeader("X-Public") == "true")

	var usersResponse users.UsersResponse

	usersResponse.Total = *total
	usersResponse.Users = users_result

	c.JSON(http.StatusOK, usersResponse)
}

func (cont *userController) GetAttendants(c *gin.Context) {
	search := c.Param("search")

	userList, getErr := services.UsersService.GetAttendants(search)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	c.JSON(http.StatusOK, userList)
}
