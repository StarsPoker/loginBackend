package menus

import (
	"net/http"
	"strconv"

	"github.com/StarsPoker/loginBackend/domain/menus"
	"github.com/StarsPoker/loginBackend/services"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
	"github.com/gin-gonic/gin"
)

var (
	MenusController MenusInterface = &menusController{}
)

type MenusInterface interface {
	InsertMenu(c *gin.Context)
	GetMenus(c *gin.Context)
	GetChildrens(c *gin.Context)
	BuildMenu(c *gin.Context)
	DeleteMenu(c *gin.Context)
	UpdateMenu(c *gin.Context)
	ChangeOrderUpMenu(c *gin.Context)
	ChangeOrderDownMenu(c *gin.Context)
}

type menusController struct {
}

func (ctrl *menusController) GetMenus(c *gin.Context) {

	menuProfileFatherParam := c.Param("profile_father")
	menuProfileFather, MenuErr := strconv.ParseInt(menuProfileFatherParam, 10, 64)

	if MenuErr != nil {
		restErr := rest_errors.NewBadRequestError("menu id should be a number (ProfileFather)")
		c.JSON(restErr.Status, restErr)
		return
	}

	var m menus.Menu
	m.ProfileFather = &menuProfileFather

	result, errGet := services.MenusService.GetMenus(m)
	if errGet != nil {
		c.JSON(errGet.Status, errGet)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (ctrl *menusController) GetChildrens(c *gin.Context) {

	menuProfileFatherParam := c.Query("profileFather")
	menuIdParam := c.Query("id")

	menuId, MenuErro := strconv.ParseInt(menuIdParam, 10, 64)
	menuProfileFather, MenuErr := strconv.ParseInt(menuProfileFatherParam, 10, 64)

	if MenuErro != nil {
		restErr := rest_errors.NewBadRequestError("menu id should be a number (id)")
		c.JSON(restErr.Status, restErr)
		return
	}

	if menuProfileFather != 0 {
		if MenuErr != nil {
			restErr := rest_errors.NewBadRequestError("menu id should be a number (ProfileFather)")
			c.JSON(restErr.Status, restErr)
			return
		}
	}

	var m menus.Menu
	m.Id = menuId
	m.ProfileFather = &menuProfileFather

	result, errGet := services.MenusService.GetChildrens(m)
	if errGet != nil {
		c.JSON(errGet.Status, errGet)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (ctrl *menusController) DeleteMenu(c *gin.Context) {
	menuIdParam := c.Param("menu_id")
	menuId, MenuErr := strconv.ParseInt(menuIdParam, 10, 64)

	if MenuErr != nil {
		restErr := rest_errors.NewBadRequestError("menu id should be a number")
		c.JSON(restErr.Status, restErr)
		return
	}

	var menuDelete menus.Menu
	menuDelete.Id = menuId

	errDelete := services.MenusService.DeleteMenu(menuDelete)
	if errDelete != nil {
		c.JSON(errDelete.Status, errDelete)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"deleted": "true"})
	return
}

func (ctrl *menusController) InsertMenu(c *gin.Context) {
	var m menus.Menu
	if err := c.ShouldBindJSON(&m); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	result, saveErr := services.MenusService.InsertMenu(m)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (ctrl *menusController) UpdateMenu(c *gin.Context) {
	menuIdParam := c.Param("menu_id")
	menuId, menuErr := strconv.ParseInt(menuIdParam, 10, 64)

	if menuErr != nil {
		restErr := rest_errors.NewBadRequestError("menu id should be a number")
		c.JSON(restErr.Status, restErr)
		return
	}

	var m menus.Menu
	if err := c.ShouldBindJSON(&m); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	m.Id = menuId

	result, updateErr := services.MenusService.UpdateMenu(m)
	if updateErr != nil {
		c.JSON(updateErr.Status, updateErr)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (ctrl *menusController) ChangeOrderUpMenu(c *gin.Context) {
	menuIdParam := c.Param("menu_id")
	menuId, menuErr := strconv.ParseInt(menuIdParam, 10, 64)

	if menuErr != nil {
		restErr := rest_errors.NewBadRequestError("menu id should be a number")
		c.JSON(restErr.Status, restErr)
		return
	}

	var m menus.Menu
	if err := c.ShouldBindJSON(&m); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	m.Id = menuId

	result, updateErr := services.MenusService.ChangeOrderUpMenu(m)
	if updateErr != nil {
		c.JSON(updateErr.Status, updateErr)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (ctrl *menusController) ChangeOrderDownMenu(c *gin.Context) {
	menuIdParam := c.Param("menu_id")
	menuId, menuErr := strconv.ParseInt(menuIdParam, 10, 64)

	if menuErr != nil {
		restErr := rest_errors.NewBadRequestError("menu id should be a number")
		c.JSON(restErr.Status, restErr)
		return
	}

	var m menus.Menu
	if err := c.ShouldBindJSON(&m); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	m.Id = menuId

	result, updateErr := services.MenusService.ChangeOrderDownMenu(m)
	if updateErr != nil {
		c.JSON(updateErr.Status, updateErr)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (cont *menusController) BuildMenu(c *gin.Context) {
	token := c.Request.Header["Authorization"][0]
	at, tokenErr := services.AccessTokenService.GetById(token)
	if tokenErr != nil {
		err := rest_errors.NewBadRequestError("Invalid access token")
		c.JSON(err.Status, err)
		return
	}

	menus, getErr := services.MenusService.BuildMenu(at.UserId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}

	c.JSON(http.StatusOK, menus)
}
