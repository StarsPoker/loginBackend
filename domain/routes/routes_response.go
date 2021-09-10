package routes

type RoutesResponse struct {
	Total  int     `json:"total"`
	Routes []Route `json:"data"`
}
