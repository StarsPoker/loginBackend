package users

type UsersResponse struct {
	Page         int64         `json:"page"`
	ItemsPerPage int64         `json:"items_per_page"`
	Users        []interface{} `json:"data"`
}
