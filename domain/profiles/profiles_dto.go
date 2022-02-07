package profiles

import "github.com/StarsPoker/loginBackend/utils/errors/rest_errors"

type Profiles []Profile

type Profile struct {
	Id               int64  `json:"id"`
	Name             string `json:"name"`
	ProfileCode      string `json:"profile_code"`
	Withdrawal       *int64 `json:"withdrawal"`
	Expense          *int64 `json:"expense"`
	Bot              *int64 `json:"bot"`
	Closure          *int64 `json:"closure"`
	Atendence        *int64 `json:"atendence"`
	FinishWithdrawal *int64 `json:"finish_withdrawal"`
}

type Filter struct {
	Id          string
	Name        string
	ProfileCode string
}

func (p *Profile) Validate() *rest_errors.RestErr {
	return nil
}

type Users []User

type User struct {
	Id        *int64 `json:"id"`
	Name      string `json:"name"`
	Role      int64  `json:"role"`
	Status    int64  `json:"status"`
	IdProfile *int64 `json:"id_profile"`
}

type Routes []Route

type Route struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Type   int64  `json:"type"`
	MenuSt string `json:"menu_string"`
	MenuId int64  `json:"menu"`
}

type ProfilesUsers []ProfileUser

type ProfileUser struct {
	Id        int64 `json:"id"`
	IdUser    int64 `json:"id_user"`
	IdProfile int64 `json:"id_profile"`
}

type ProfilesRoutes []ProfileRoute

type ProfileRoute struct {
	Id        int64 `json:"id"`
	IdRoute   int64 `json:"id_route"`
	IdProfile int64 `json:"id_profile"`
}

func (pu *ProfileUser) Validate() *rest_errors.RestErr {
	// if pu.IdProfile != -1 {
	// 	return rest_errors.NewBadRequestError("invalid access token id")
	// }

	return nil
}

func (pu *ProfileRoute) Validate() *rest_errors.RestErr {
	// if pu.IdProfile != -1 {
	// 	return rest_errors.NewBadRequestError("invalid access token id")
	// }

	return nil
}

type ProfilesMenus []ProfileMenu

type ProfileMenu struct {
	Id        int64 `json:"id"`
	IdMenu    int64 `json:"id_menu"`
	IdProfile int64 `json:"id_profile"`
}

func (pm *ProfileMenu) Validate() *rest_errors.RestErr {
	// if pu.IdProfile != -1 {
	// 	return rest_errors.NewBadRequestError("invalid access token id")
	// }

	return nil
}

type BuildMenus []BuildMenu

type BuildMenu struct {
	Id          *int64      `json:"id"`
	Icon        *string     `json:"icon"`
	Parent      *int64      `json:"parent"`
	Description string      `json:"description"`
	Level       int64       `json:"level"`
	Link        *string     `json:"link"`
	HasSubGroup *bool       `json:"hasSubGroup"`
	Menus       []BuildMenu `json:"menus"`
}
