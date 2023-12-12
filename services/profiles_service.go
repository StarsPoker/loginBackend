package services

import (
	"fmt"

	"github.com/StarsPoker/loginBackend/domain/menus"
	"github.com/StarsPoker/loginBackend/domain/profiles"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
)

var (
	ProfilesService profilesInterface = &profilesService{}
)

type profilesService struct {
}

type profilesInterface interface {
	GetProfile(int64) (*profiles.Profile, *rest_errors.RestErr)
	GetProfiles(int, int, *profiles.Filter, int64) (profiles.Profiles, *int, *rest_errors.RestErr)
	GetProfileUser(int64) (*profiles.ProfileUser, *rest_errors.RestErr)
	GetProfileMenu(int64) (*profiles.ProfileMenu, *rest_errors.RestErr)
	GetProfileUsers(int, int, *profiles.Filter, int64) (profiles.Users, *int, *rest_errors.RestErr)
	GetProfileRoutes(int, int, *profiles.Filter, int64) (profiles.Routes, *int, *rest_errors.RestErr)
	GetProfileUsersAdds(int, int, *profiles.Filter, int64) (profiles.Users, *rest_errors.RestErr)
	GetProfilePermissions(string) (*profiles.Profile, *rest_errors.RestErr)
	GetProfileAttendants(search string, profileId int64) (profiles.Users, *rest_errors.RestErr)
	GetProfileRoutesAdds(search string, profileId int64) (profiles.Routes, *rest_errors.RestErr)
	CreateProfile(profiles.Profile) (*profiles.Profile, *rest_errors.RestErr)
	CreateProfileUser(profiles.ProfileUser) (*profiles.ProfileUser, *rest_errors.RestErr)
	CreateProfileRoute(profiles.ProfileRoute) (*profiles.ProfileRoute, *rest_errors.RestErr)
	CreateProfileMenu(profiles.ProfileMenu) (*profiles.ProfileMenu, *rest_errors.RestErr)
	CreateProfileMenuFather(profiles.ProfileMenu) (*profiles.ProfileMenu, *rest_errors.RestErr)
	UpdateProfileUser(profiles.ProfileUser) (*profiles.ProfileUser, *rest_errors.RestErr)
	UpdateProfile(profiles.Profile) (*profiles.Profile, *rest_errors.RestErr)
	UpdateParam(profiles.Profile) (*profiles.Profile, *rest_errors.RestErr)
	DeleteProfile(p profiles.Profile) *rest_errors.RestErr
	DeleteProfileUser(p profiles.ProfileUser) *rest_errors.RestErr
	DeleteProfileRoute(p profiles.ProfileRoute) *rest_errors.RestErr
	DeleteProfileMenu(p profiles.ProfileMenu) *rest_errors.RestErr
	DeleteProfileMenuFather(p profiles.ProfileMenu) *rest_errors.RestErr
}

func (s *profilesService) GetProfile(profileId int64) (*profiles.Profile, *rest_errors.RestErr) {
	result := &profiles.Profile{Id: profileId}
	if err := result.GetProfile(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *profilesService) GetProfileUser(profileUserId int64) (*profiles.ProfileUser, *rest_errors.RestErr) {
	result := &profiles.ProfileUser{IdUser: profileUserId}
	if err := result.GetProfileUser(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *profilesService) GetProfileUser2(profileUserId int64) (*profiles.ProfileUser, *rest_errors.RestErr) {
	result := &profiles.ProfileUser{IdUser: profileUserId}
	if err := result.GetProfileUser2(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *profilesService) GetProfileRoute(profileRouteId int64) (*profiles.ProfileRoute, *rest_errors.RestErr) {
	result := &profiles.ProfileRoute{IdRoute: profileRouteId}
	if err := result.GetProfileRoute(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *profilesService) GetProfileMenu(profileMenuId int64) (*profiles.ProfileMenu, *rest_errors.RestErr) {
	result := &profiles.ProfileMenu{IdMenu: profileMenuId}
	if err := result.GetProfileMenu(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *profilesService) GetProfiles(page int, itemsPerPage int, filter *profiles.Filter, userId int64) (profiles.Profiles, *int, *rest_errors.RestErr) {
	dao := &profiles.Profile{}
	profiles, total, err := dao.GetProfiles(page, itemsPerPage, filter, userId)
	if err != nil {
		return nil, nil, err
	}

	return profiles, total, nil
}

func (s *profilesService) CreateProfile(profileToSave profiles.Profile) (*profiles.Profile, *rest_errors.RestErr) {

	if err := profileToSave.Validate(); err != nil {
		return nil, err
	}

	if err := profileToSave.Save(); err != nil {
		return nil, err
	}
	return &profileToSave, nil
}

func (s *profilesService) UpdateProfile(p profiles.Profile) (*profiles.Profile, *rest_errors.RestErr) {
	current, err := s.GetProfile(p.Id)
	if err != nil {
		return nil, err
	}

	current.ProfileCode = p.ProfileCode
	current.Name = p.Name
	current.Id = p.Id

	if err := current.Update(); err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *profilesService) UpdateParam(p profiles.Profile) (*profiles.Profile, *rest_errors.RestErr) {
	current, err := s.GetProfile(p.Id)
	if err != nil {
		return nil, err
	}

	if p.Withdrawal != nil {
		current.Withdrawal = p.Withdrawal
	} else if p.Expense != nil {
		current.Expense = p.Expense
	} else if p.Bot != nil {
		current.Bot = p.Bot
	} else if p.Atendence != nil {
		current.Atendence = p.Atendence
	} else if p.Closure != nil {
		current.Closure = p.Closure
	} else if p.FinishWithdrawal != nil {
		current.FinishWithdrawal = p.FinishWithdrawal
	} else if p.MakeBlockedWithdrawal != nil {
		current.MakeBlockedWithdrawal = p.MakeBlockedWithdrawal
	} else if p.MakeAlertWithdrawal != nil {
		current.MakeAlertWithdrawal = p.MakeAlertWithdrawal
	}else if p.FinishAtendance != nil {
		current.FinishAtendance = p.FinishAtendance
	}

	if err := current.UpdateParam(); err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *profilesService) DeleteProfile(p profiles.Profile) *rest_errors.RestErr {
	current, err := s.GetProfile(p.Id)
	if err != nil {
		return err
	}

	if err := current.Delete(); err != nil {
		return nil
	}

	return nil
}

func (s *profilesService) GetProfileUsers(page int, itemsPerPage int, filter *profiles.Filter, profileId int64) (profiles.Users, *int, *rest_errors.RestErr) {
	dao := &profiles.Profile{}
	profiles, total, err := dao.GetProfileUsers(page, itemsPerPage, filter, profileId)
	if err != nil {
		return nil, nil, err
	}

	return profiles, total, nil
}

func (s *profilesService) GetProfileRoutes(page int, itemsPerPage int, filter *profiles.Filter, profileId int64) (profiles.Routes, *int, *rest_errors.RestErr) {
	dao := &profiles.Profile{}
	profiles, total, err := dao.GetProfileRoutes(page, itemsPerPage, filter, profileId)
	if err != nil {
		return nil, nil, err
	}

	return profiles, total, nil
}

func (s *profilesService) GetProfileUsersAdds(page int, itemsPerPage int, filter *profiles.Filter, profileId int64) (profiles.Users, *rest_errors.RestErr) {
	dao := &profiles.Profile{}
	profiles, err := dao.GetProfileUsersAdds(page, itemsPerPage, filter, profileId)
	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func (s *profilesService) GetProfileAttendants(search string, profileId int64) (profiles.Users, *rest_errors.RestErr) {
	dao := &profiles.User{}
	profiles, err := dao.GetProfileAttendants(search, profileId)
	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func (s *profilesService) GetProfileRoutesAdds(search string, profileId int64) (profiles.Routes, *rest_errors.RestErr) {
	dao := &profiles.Route{}
	profiles, err := dao.GetProfileRoutesAdds(search, profileId)
	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func (s *profilesService) GetProfilePermissions(profileId string) (*profiles.Profile, *rest_errors.RestErr) {

	result := &profiles.Profile{ProfileCode: profileId}
	if err := result.GetProfilePermissions(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *profilesService) CreateProfileUser(profileToSave profiles.ProfileUser) (*profiles.ProfileUser, *rest_errors.RestErr) {

	if err := profileToSave.Validate(); err != nil {
		return nil, err
	}

	if err := profileToSave.SaveProfileUser(); err != nil {
		return nil, err
	}
	return &profileToSave, nil
}

func (s *profilesService) CreateProfileRoute(profileToSave profiles.ProfileRoute) (*profiles.ProfileRoute, *rest_errors.RestErr) {

	if err := profileToSave.Validate(); err != nil {
		return nil, err
	}

	if err := profileToSave.SaveProfileRoute(); err != nil {
		return nil, err
	}
	return &profileToSave, nil
}

func (s *profilesService) CreateProfileMenu(profileToSave profiles.ProfileMenu) (*profiles.ProfileMenu, *rest_errors.RestErr) {
	if err := profileToSave.Validate(); err != nil {
		return nil, err
	}

	result := &menus.Menu{Id: profileToSave.IdMenu}
	if err := result.GetMenu(); err != nil {
		return nil, err
	}

	fatherId := *result.Parent
	total, err := profileToSave.GetTotalProfileMenu(fatherId)
	if err != nil {
		return nil, err
	}

	if *total == 0 {
		profileFather := &profiles.ProfileMenu{IdMenu: fatherId, IdProfile: profileToSave.IdProfile}
		if err := profileFather.SaveProfileMenu(); err != nil {
			return nil, err
		}
	}

	if err := profileToSave.SaveProfileMenu(); err != nil {
		return nil, err
	}

	if err := profileToSave.DeleteRoutesRelation(); err != nil {
		return nil, err
	}

	if err := profileToSave.SaveRoutesRelation(); err != nil {
		return nil, err
	}

	return &profileToSave, nil
}

func (s *profilesService) CreateProfileMenuFather(profileToSave profiles.ProfileMenu) (*profiles.ProfileMenu, *rest_errors.RestErr) {
	if err := profileToSave.Validate(); err != nil {
		return nil, err
	}

	if err := profileToSave.SaveProfileMenu(); err != nil {
		return nil, err
	}
	return &profileToSave, nil
}

func (s *profilesService) UpdateProfileUser(p profiles.ProfileUser) (*profiles.ProfileUser, *rest_errors.RestErr) {
	current, err := s.GetProfileUser(p.Id)
	if err != nil {
		return nil, err
	}

	current.IdProfile = p.IdProfile
	current.IdUser = p.IdUser

	if err := current.UpdateProfileUser(); err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *profilesService) DeleteProfileUser(p profiles.ProfileUser) *rest_errors.RestErr {
	current, err := s.GetProfileUser2(p.Id)
	if err != nil {
		return err
	}

	if err := current.DeleteProfileUser(); err != nil {
		return nil
	}

	return nil
}

func (s *profilesService) DeleteProfileRoute(p profiles.ProfileRoute) *rest_errors.RestErr {
	current, err := s.GetProfileRoute(p.Id)
	if err != nil {
		return err
	}

	if err := current.DeleteProfileRoute(); err != nil {
		return nil
	}

	return nil
}

func (s *profilesService) DeleteProfileMenu(p profiles.ProfileMenu) *rest_errors.RestErr {
	current, err := s.GetProfileMenu(p.Id)
	if err != nil {
		return err
	}

	result := &menus.Menu{Id: current.IdMenu}
	if err := result.GetMenu(); err != nil {
		return err
	}

	var busca menus.Menu
	busca.ProfileFather = &current.IdProfile
	busca.Id = *result.Parent

	brothers, err := busca.GetChildrens()
	if err != nil {
		return err
	}

	var count = 0
	for i, s := range brothers {
		if s.HasRelation == 1 {
			count = count + 1
		}
		fmt.Println(i)
	}

	var buscaFather profiles.ProfileMenu
	buscaFather.IdProfile = current.IdProfile
	buscaFather.IdMenu = *result.Parent

	if count == 1 {
		father, err := buscaFather.GetProfileMenuFather()
		if err != nil {
			return err
		}

		var fatherDelete = father[0]

		if err := fatherDelete.DeleteProfileMenu(); err != nil {
			return nil
		}
	}

	if err := current.DeleteProfileMenu(); err != nil {
		return nil
	}

	if err := current.DeleteRoutesRelation(); err != nil {
		return nil
	}

	return nil
}

func (s *profilesService) DeleteProfileMenuFather(p profiles.ProfileMenu) *rest_errors.RestErr {
	current, err := s.GetProfileMenu(p.Id)
	if err != nil {
		return err
	}

	if err := current.DeleteProfileMenu(); err != nil {
		return nil
	}

	return nil
}
