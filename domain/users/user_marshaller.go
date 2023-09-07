package users

type PublicUser struct {
	Name          string  `json:"name"`
	Role          int64   `json:"role"`
	Status        int64   `json:"status"`
	Email         string  `json:"email"`
	ProfileAccess *string `json:"profile_access"`
	Inscription   *string `json:"inscription"`
}

type PrivateUser struct {
	AuthenticatorConfigured bool    `json:"authenticator_configured"`
	Id                      int64   `json:"id"`
	Name                    string  `json:"name"`
	Email                   string  `json:"email"`
	InstanceId              *string `json:"instance_id"`
	InstanceName            *string `json:"instance_name"`
	Contact                 *string `json:"contact"`
	Role                    int64   `json:"role"`
	Status                  int64   `json:"status"`
	DateCreated             string  `json:"date_created"`
	DefaultPassword         int64   `json:"default_password"`
	ProfileAccess           *string `json:"profile_access"`
	Inscription             *string `json:"inscription"`
}

func (users Users) Marshall(isPublic bool) []interface{} {
	result := make([]interface{}, len(users))

	for index, user := range users {
		result[index] = user.Marshall(isPublic)
	}
	return result
}

func (user *User) Marshall(isPublic bool) interface{} {

	role := int64(0)
	if user.Role != nil {
		role = *user.Role
	}

	if isPublic {
		return PublicUser{
			Name:   user.Name,
			Role:   role,
			Status: user.Status,
		}
	}

	return PrivateUser{
		Id:                      user.Id,
		Name:                    user.Name,
		Email:                   user.Email,
		Role:                    role,
		Status:                  user.Status,
		Contact:                 user.Contact,
		InstanceId:              user.InstanceId,
		InstanceName:            user.InstanceName,
		DateCreated:             user.DateCreated,
		DefaultPassword:         user.DefaultPassword,
		ProfileAccess:           user.ProfileAccess,
		Inscription:             user.Inscription,
		AuthenticatorConfigured: user.AuthenticatorConfigured,
	}
}
