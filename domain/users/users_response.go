package users

type UsersResponse struct {
	Total int           `json:"total"`
	Users []interface{} `json:"data"`
}
