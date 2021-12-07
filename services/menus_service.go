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
	for i, profile := range profileRelation {
		if profile.Parent == nil {
			hs := false
			profile.HasSubGroup = &hs
			buildMenus = append(buildMenus, profile)
			i++
		} else {
			for a, father := range buildMenus {
				if *father.Id == *profile.Parent {
					hs := true
					buildMenus[a].HasSubGroup = &hs
					buildMenus[a].Menus = append(buildMenus[a].Menus, profile)
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
