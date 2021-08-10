package menus

import "github.com/StarsPoker/loginBackend/utils/errors/rest_errors"

type Menus []Menu

type Menu struct {
	Id            int64   `json:"id"`
	Icon          *string `json:"icon"`
	Parent        *int64  `json:"parent"`
	Name          string  `json:"name"`
	Link          *string `json:"link"`
	Level         int64   `json:"level"`
	Order         int64   `json:"order"`
	HasChield     int64   `json:"has_chield"`
	HasRelation   int64   `json:"has_relation"`
	ProfileFather *int64  `json:"profile_father"`
	IdRelation    int64   `json:"id_relation"`
}

type MaxOrder struct {
	MaxOrder int64 `json:"max_order"`
}

func (m *Menu) Validate() *rest_errors.RestErr {
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
