package services

import (
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
	GetProfiles(int, int, *profiles.Filter) (profiles.Profiles, *int, *rest_errors.RestErr)
	GetProfileUser(int64) (*profiles.ProfileUser, *rest_errors.RestErr)
	GetProfileMenu(int64) (*profiles.ProfileMenu, *rest_errors.RestErr)
	GetProfileUsers(int, int, *profiles.Filter, int64) (profiles.Users, *int, *rest_errors.RestErr)
	GetProfileUsersAdds(int, int, *profiles.Filter, int64) (profiles.Users, *rest_errors.RestErr)
	GetProfileAttendants(search string, profileId int64) (profiles.Users, *rest_errors.RestErr)
	CreateProfile(profiles.Profile) (*profiles.Profile, *rest_errors.RestErr)
	CreateProfileUser(profiles.ProfileUser) (*profiles.ProfileUser, *rest_errors.RestErr)
	CreateProfileMenu(profiles.ProfileMenu) (*profiles.ProfileMenu, *rest_errors.RestErr)
	CreateProfileMenuFather(profiles.ProfileMenu) (*profiles.ProfileMenu, *rest_errors.RestErr)
	UpdateProfileUser(profiles.ProfileUser) (*profiles.ProfileUser, *rest_errors.RestErr)
	UpdateProfile(profiles.Profile) (*profiles.Profile, *rest_errors.RestErr)
	DeleteProfile(p profiles.Profile) *rest_errors.RestErr
	DeleteProfileUser(p profiles.ProfileUser) *rest_errors.RestErr
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

func (s *profilesService) GetProfileMenu(profileMenuId int64) (*profiles.ProfileMenu, *rest_errors.RestErr) {
	result := &profiles.ProfileMenu{IdMenu: profileMenuId}
	if err := result.GetProfileMenu(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *profilesService) GetProfiles(page int, itemsPerPage int, filter *profiles.Filter) (profiles.Profiles, *int, *rest_errors.RestErr) {
	dao := &profiles.Profile{}
	profiles, total, err := dao.GetProfiles(page, itemsPerPage, filter)
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

func (s *profilesService) CreateProfileUser(profileToSave profiles.ProfileUser) (*profiles.ProfileUser, *rest_errors.RestErr) {

	if err := profileToSave.Validate(); err != nil {
		return nil, err
	}

	if err := profileToSave.SaveProfileUser(); err != nil {
		return nil, err
	}
	return &profileToSave, nil
}

func (s *profilesService) CreateProfileMenu(profileToSave profiles.ProfileMenu) (*profiles.ProfileMenu, *rest_errors.RestErr) {
	if err := profileToSave.Validate(); err != nil {
		return nil, err
	}

	// Buscar o menu que vou relacionar
	result := &menus.Menu{Id: profileToSave.IdMenu}
	if err := result.GetMenu(); err != nil {
		return nil, err
	}

	// Buscar o pai do menu que vou relacionar
	fatherId := *result.Parent

	// Ver se o pai do menu que vou relacionar está relacionado ou nao
	total, err := profileToSave.GetTotalProfileMenu(fatherId)
	if err != nil {
		return nil, err
	}

	// if (não esta relacionado) {
	//	relacionar o id do pai com o id do profile que estou recebendo
	// }

	if *total == 0 {
		profileFather := &profiles.ProfileMenu{IdMenu: fatherId, IdProfile: profileToSave.IdProfile}
		if err := profileFather.SaveProfileMenu(); err != nil {
			return nil, err
		}
	}

	// Depois Continuar o process normal

	if err := profileToSave.SaveProfileMenu(); err != nil {
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
	current, err := s.GetProfileUser(p.Id)
	if err != nil {
		return err
	}

	if err := current.DeleteProfileUser(); err != nil {
		return nil
	}

	return nil
}

func (s *profilesService) DeleteProfileMenu(p profiles.ProfileMenu) *rest_errors.RestErr {
	current, err := s.GetProfileMenu(p.Id)
	if err != nil {
		return err
	}

	// Buscar o menu que vou apagar
	result := &menus.Menu{Id: current.IdMenu}
	if err := result.GetMenu(); err != nil {
		return err
	}

	var busca menus.Menu
	busca.ProfileFather = &current.IdProfile
	busca.Id = *result.Parent

	// Buscar os menus irmaos
	brothers, err := busca.GetChildrens()
	if err != nil {
		return err
	}

	// ver se existe mais de um filho relacionado
	var count = 0
	for i, s := range brothers {
		if s.HasRelation == 1 {
			count = count + i
		}
	}

	// if (nao existir irmaos relacionados) {
	//	apagar o id do pai com o id do profile que estou recebendo
	// }
	var buscaFather profiles.ProfileMenu
	buscaFather.IdProfile = current.IdProfile
	buscaFather.IdMenu = *result.Parent

	if count == 0 {
		father, err := buscaFather.GetProfileMenuFather()
		if err != nil {
			return err
		}

		var fatherDelete = father[0]

		if err := fatherDelete.DeleteProfileMenu(); err != nil {
			return nil
		}
	}

	// Depois Continuar o process normal
	if err := current.DeleteProfileMenu(); err != nil {
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
