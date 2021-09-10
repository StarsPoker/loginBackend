package profiles

type ProfilesResponse struct {
	Total    int       `json:"total"`
	Profiles []Profile `json:"data"`
}

type UsersResponse struct {
	Total int    `json:"total"`
	Users []User `json:"data"`
}

type ProfilesUsersResponse struct {
	Total         int           `json:"total"`
	ProfilesUsers []ProfileUser `json:"data"`
}

type RoutesResponse struct {
	Total  int     `json:"total"`
	Routes []Route `json:"data"`
}
