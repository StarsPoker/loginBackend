package routes

import "github.com/StarsPoker/loginBackend/utils/errors/rest_errors"

type Routes []Route

type Route struct {
	Id     int64   `json:"id"`
	Name   string  `json:"name"`
	Type   int64   `json:"type"`
	MenuSt *string `json:"menu_string"`
	MenuId int64   `json:"menu"`
}

type Filter struct {
	Id     string
	Name   string
	Type   string
	MenuId string
}

func (r *Route) Validate() *rest_errors.RestErr {
	return nil
}
