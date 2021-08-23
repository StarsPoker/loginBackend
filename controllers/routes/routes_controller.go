package routes

import (
	"net/http"
	"strconv"

	"github.com/StarsPoker/loginBackend/domain/routes"
	"github.com/StarsPoker/loginBackend/services"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
	"github.com/gin-gonic/gin"
)

var (
	RouteController RouteInterface = &routeController{}
)

type RouteInterface interface {
	GetRoutes(c *gin.Context)
	CreateRoute(c *gin.Context)
	UpdateRoute(c *gin.Context)
	DeleteRoute(c *gin.Context)
}

type routeController struct {
}

func (cont *routeController) GetRoutes(c *gin.Context) {

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

	result, total, getErr := services.RoutesService.GetRoutes(page, itemsPerPage, filter)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}

	var routesResponse routes.RoutesResponse

	routesResponse.Total = *total
	routesResponse.Routes = result

	c.JSON(http.StatusOK, routesResponse)
}

func buildFilter(c *gin.Context) *routes.Filter {
	var filter routes.Filter
	filter.Id = c.Query("id")
	filter.Name = c.Query("name")
	filter.Type = c.Query("type")
	filter.MenuId = c.Query("menu")
	return &filter
}

func (cont *routeController) CreateRoute(c *gin.Context) {
	var w routes.Route

	if err := c.ShouldBindJSON(&w); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	result, saveErr := services.RoutesService.CreateRoute(w)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (cont *routeController) UpdateRoute(c *gin.Context) {
	routeIdParam := c.Param("route_id")
	routeId, routeErr := strconv.ParseInt(routeIdParam, 10, 64)

	if routeErr != nil {
		restErr := rest_errors.NewBadRequestError("route id should be a number")
		c.JSON(restErr.Status, restErr)
		return
	}

	var b routes.Route
	if err := c.ShouldBindJSON(&b); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	b.Id = routeId

	result, updateErr := services.RoutesService.UpdateRoute(b)
	if updateErr != nil {
		c.JSON(updateErr.Status, updateErr)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (cont *routeController) DeleteRoute(c *gin.Context) {
	routeIdParam := c.Param("route_id")

	routeId, routeErr := strconv.ParseInt(routeIdParam, 10, 64)

	if routeErr != nil {
		restErr := rest_errors.NewBadRequestError("route id should be a number")
		c.JSON(restErr.Status, restErr)
		return
	}

	var routeDelete routes.Route

	routeDelete.Id = routeId

	deleteErr := services.RoutesService.DeleteRoute(routeDelete)
	if deleteErr != nil {
		c.JSON(deleteErr.Status, deleteErr)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}
