package services

import (
	"github.com/StarsPoker/loginBackend/domain/menus"
	"github.com/StarsPoker/loginBackend/domain/profiles"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
)

var (
	MenusService menusInterface = &menusService{}
)

type menusService struct {
}

type menusInterface interface {
	InsertMenu(menus.Menu) (*menus.Menu, *rest_errors.RestErr)
	GetChildrens(menus.Menu) (menus.Menus, *rest_errors.RestErr)
	GetMenus(menus.Menu) (menus.Menus, *rest_errors.RestErr)
	GetMenu(int64) (*menus.Menu, *rest_errors.RestErr)
	UpdateMenu(menus.Menu) (*menus.Menu, *rest_errors.RestErr)
	ChangeOrderUpMenu(menus.Menu) (*menus.Menu, *rest_errors.RestErr)
	ChangeOrderDownMenu(menus.Menu) (*menus.Menu, *rest_errors.RestErr)
	DeleteMenu(menus.Menu) *rest_errors.RestErr
	BuildMenu(int64) ([]profiles.BuildMenu, *rest_errors.RestErr)
	BuildMenuSearch(int64, string) ([]profiles.BuildMenu, *rest_errors.RestErr)
	ProfilePermission(int64, string) (*menus.Permission, *rest_errors.RestErr)
	GetChildrenSearch(search string) (menus.Menus, *rest_errors.RestErr)
	GetProfilesRelation(menuId int64) ([]menus.ProfileRelation, *rest_errors.RestErr)
}

func (s *menusService) BuildMenu(acessToken int64) ([]profiles.BuildMenu, *rest_errors.RestErr) {

	profileBusca := &profiles.Profile{Id: acessToken}

	profileRelation, err := profileBusca.GetProfileRelation()
	if err != nil {
		return nil, err
	}

	buildMenus := make([]profiles.BuildMenu, 0)
	for _, profile := range profileRelation {
		if profile.Parent == nil {
			hs := false
			profile.HasSubGroup = &hs
			buildMenus = append(buildMenus, profile)
		} else if profile.Level == 2 {
			for a, father := range buildMenus {
				if *father.Id == *profile.Parent {
					hs := true
					buildMenus[a].HasSubGroup = &hs
					buildMenus[a].Menus = append(buildMenus[a].Menus, profile)
				}
			}
		} else if profile.Level == 3 {
			for i, menu := range buildMenus {
				if len(menu.Menus) > 0 {
					for j, subMenu := range menu.Menus {
						if *subMenu.Id == *profile.Parent {
							verd := true
							buildMenus[i].Menus[j].HasSubGroup = &verd
							buildMenus[i].Menus[j].Menus = append(buildMenus[i].Menus[j].Menus, profile)
						}
					}
				}
			}
		}
	}

	return buildMenus, nil
}

func (s *menusService) BuildMenuSearch(acessToken int64, menuSearch string) ([]profiles.BuildMenu, *rest_errors.RestErr) {

	profileBusca := &profiles.Profile{Id: acessToken}

	profileRelation, err := profileBusca.GetProfileRelationSearch(menuSearch)
	if err != nil {
		return nil, err
	}

	buildMenus := make([]profiles.BuildMenu, 0)
	for _, profile := range profileRelation {
		if profile.Level == 2 {
			hasFather := 0
			for a, father := range buildMenus {
				if *father.Id == *profile.Parent {
					hs := true
					buildMenus[a].HasSubGroup = &hs
					buildMenus[a].Menus = append(buildMenus[a].Menus, profile)
					hasFather = 1
				}
			}
			if hasFather == 0 {
				menuFather, menuFatherErr := profileBusca.GetMenuFather(*profile.Parent)
				if menuFatherErr != nil {
					return nil, menuFatherErr
				}

				hs := false
				menuFather[0].HasSubGroup = &hs
				buildMenus = append(buildMenus, menuFather[0])
				for a, father := range buildMenus {
					if *father.Id == *profile.Parent {
						hs := true
						buildMenus[a].HasSubGroup = &hs
						buildMenus[a].Menus = append(buildMenus[a].Menus, profile)
						hasFather = 1
					}
				}
			}
		} else if profile.Level == 3 {
			itsok := 0
			for a, grandFather := range buildMenus {
				if len(grandFather.Menus) > 0 {
					for b, father := range grandFather.Menus {
						if *father.Id == *profile.Parent {
							hs := true
							buildMenus[a].HasSubGroup = &hs
							buildMenus[a].Menus[b].HasSubGroup = &hs
							buildMenus[a].Menus[b].Menus = append(buildMenus[a].Menus[b].Menus, profile)
							itsok = 1
						}
					}
				}
			}
			if itsok == 0 {
				itsok2 := 0
				menuFather, menuFatherErr := profileBusca.GetMenuFather(*profile.Parent)
				if menuFatherErr != nil {
					return nil, menuFatherErr
				}

				for i, father := range buildMenus {
					if *father.Id == *menuFather[0].Parent {
						hs := true
						buildMenus[i].HasSubGroup = &hs
						menuFather[0].HasSubGroup = &hs
						menuFather[0].Menus = append(menuFather[0].Menus, profile)
						buildMenus[i].Menus = append(buildMenus[i].Menus, menuFather[0])
						itsok2 = 1
					}
				}

				if itsok2 == 0 {
					menuGrandFather, menuGrandFatherErr := profileBusca.GetMenuFather(*menuFather[0].Parent)
					if menuGrandFatherErr != nil {
						return nil, menuGrandFatherErr
					}
					hs := true
					menuGrandFather[0].HasSubGroup = &hs
					menuFather[0].HasSubGroup = &hs
					menuFather[0].Menus = append(menuFather[0].Menus, profile)
					menuGrandFather[0].Menus = append(menuGrandFather[0].Menus, menuFather[0])
					buildMenus = append(buildMenus, menuGrandFather[0])
				}
			}
		}
	}

	return buildMenus, nil
}

func (s *menusService) ProfilePermission(acessToken int64, menuName string) (*menus.Permission, *rest_errors.RestErr) {

	profilePermission := &menus.ProfilePermission{MenuName: menuName, UserId: acessToken}

	count, err := profilePermission.GetUserPermission()
	if err != nil {
		return nil, err
	}

	var permission menus.Permission
	permission.Permission = *count

	return &permission, nil
}

func (s *menusService) GetMenu(menuId int64) (*menus.Menu, *rest_errors.RestErr) {
	result := &menus.Menu{Id: menuId}
	if err := result.GetMenu(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *menusService) GetMenus(m menus.Menu) (menus.Menus, *rest_errors.RestErr) {

	menus, err := m.GetMenus()
	if err != nil {
		return nil, err
	}

	return menus, nil
}

func (s *menusService) GetChildrens(m menus.Menu) (menus.Menus, *rest_errors.RestErr) {

	menus, err := m.GetChildrens()
	if err != nil {
		return nil, err
	}

	return menus, nil
}

func (s *menusService) InsertMenu(menusToSave menus.Menu) (*menus.Menu, *rest_errors.RestErr) {

	if err := menusToSave.Validate(); err != nil {
		return nil, err
	}

	maxOrder := menusToSave.GetMaxOrder()

	if maxOrder == nil {
		return nil, rest_errors.NewInternalServerError("database error")
	}

	menusToSave.Order = *maxOrder

	if err := menusToSave.Save(); err != nil {
		return nil, err
	}
	return &menusToSave, nil
}

func (s *menusService) DeleteMenu(m menus.Menu) *rest_errors.RestErr {
	current, err := s.GetMenu(m.Id)
	if err != nil {
		return err
	}

	if current.Level == 1 {
		if err := current.UpdateOrder(); err != nil {
			return nil
		}

		if err := current.DeleteChildren(); err != nil {
			return nil
		}
	}

	if err := current.Delete(); err != nil {
		return nil
	}

	return nil
}

func (s *menusService) UpdateMenu(m menus.Menu) (*menus.Menu, *rest_errors.RestErr) {
	current, err := s.GetMenu(m.Id)
	if err != nil {
		return nil, err
	}

	current.Icon = m.Icon
	current.Name = m.Name
	current.Link = m.Link
	current.Id = m.Id

	if err := current.Update(); err != nil {
		return nil, err
	}

	return &m, nil
}

func (s *menusService) ChangeOrderUpMenu(m menus.Menu) (*menus.Menu, *rest_errors.RestErr) {
	current, err := s.GetMenu(m.Id)
	if err != nil {
		return nil, err
	}

	current.Order = m.Order
	current.Id = m.Id

	if err := current.UpdateOrderUpFirst(); err != nil {
		return nil, err
	}

	if err := current.UpdateOrderUpNext(); err != nil {
		return nil, err
	}

	return &m, nil
}

func (s *menusService) ChangeOrderDownMenu(m menus.Menu) (*menus.Menu, *rest_errors.RestErr) {
	current, err := s.GetMenu(m.Id)
	if err != nil {
		return nil, err
	}

	current.Order = m.Order
	current.Id = m.Id

	if err := current.UpdateOrderDownFirst(); err != nil {
		return nil, err
	}

	if err := current.UpdateOrderDownNext(); err != nil {
		return nil, err
	}

	return &m, nil
}

func (s *menusService) GetChildrenSearch(search string) (menus.Menus, *rest_errors.RestErr) {
	dao := &menus.Menus{}
	menus, err := dao.GetChildrenSearch(search)
	if err != nil {
		return nil, err
	}

	return menus, nil
}

func (s *menusService) GetProfilesRelation(menuId int64) ([]menus.ProfileRelation, *rest_errors.RestErr) {
	dao := &menus.Menus{}
	menus, err := dao.GetProfilesRelation(menuId)
	if err != nil {
		return nil, err
	}

	return menus, nil
}
