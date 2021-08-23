package services

import (
	"github.com/StarsPoker/loginBackend/domain/routes"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
)

var (
	RoutesService routesInterface = &routesService{}
)

type routesService struct {
}

type routesInterface interface {
	GetRoute(int64) (*routes.Route, *rest_errors.RestErr)
	GetRoutes(int, int, *routes.Filter) (routes.Routes, *int, *rest_errors.RestErr)
	CreateRoute(routes.Route) (*routes.Route, *rest_errors.RestErr)
	UpdateRoute(routes.Route) (*routes.Route, *rest_errors.RestErr)
	DeleteRoute(r routes.Route) *rest_errors.RestErr
}

func (s *routesService) GetRoute(routeId int64) (*routes.Route, *rest_errors.RestErr) {
	result := &routes.Route{Id: routeId}
	if err := result.GetRoute(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *routesService) GetRoutes(page int, itemsPerPage int, filter *routes.Filter) (routes.Routes, *int, *rest_errors.RestErr) {
	dao := &routes.Route{}
	routes, total, err := dao.GetRoutes(page, itemsPerPage, filter)
	if err != nil {
		return nil, nil, err
	}

	return routes, total, nil
}

func (s *routesService) CreateRoute(routeToSave routes.Route) (*routes.Route, *rest_errors.RestErr) {

	if err := routeToSave.Validate(); err != nil {
		return nil, err
	}

	if err := routeToSave.Save(); err != nil {
		return nil, err
	}
	return &routeToSave, nil
}

func (s *routesService) UpdateRoute(r routes.Route) (*routes.Route, *rest_errors.RestErr) {
	current, err := s.GetRoute(r.Id)
	if err != nil {
		return nil, err
	}

	current.MenuId = r.MenuId
	current.Type = r.Type
	current.Name = r.Name
	current.Id = r.Id

	if err := current.Update(); err != nil {
		return nil, err
	}

	return &r, nil
}

func (s *routesService) DeleteRoute(r routes.Route) *rest_errors.RestErr {
	current, err := s.GetRoute(r.Id)
	if err != nil {
		return err
	}

	if err := current.Delete(); err != nil {
		return nil
	}

	return nil
}
