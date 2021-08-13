package profiles

import (
	"github.com/StarsPoker/loginBackend/logger"

	"github.com/StarsPoker/loginBackend/datasources/mysql/stars_mysql"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
	"github.com/StarsPoker/loginBackend/utils/mysql_utils"
)

const (
	errorNoRows               = "no rows in result set"
	queryDeleteProfile        = "DELETE from profiles WHERE id = ?"
	queryGetProfile           = "SELECT id, name, profile_code FROM profiles WHERE id = ?"
	queryGetProfiles          = "SELECT p.id, p.name, p.profile_code FROM profiles p LEFT JOIN users u ON u.id = ? WHERE 1 = 1 AND p.profile_code <= u.role"
	queryInsertProfile        = "INSERT INTO profiles (name, profile_code) VALUES (?, ?)"
	queryUpdateProfile        = "UPDATE profiles SET name = ?, profile_code = ? WHERE id = ?"
	queryTotalProfiles        = "SELECT COUNT(*) as TOTAL FROM profiles p LEFT JOIN users u ON u.id = ? WHERE 1 = 1 AND p.profile_code <= u.role"
	queryGetProfileUsers      = "SELECT p.id, u.name, u.role, u.status, p.id_profile FROM users u JOIN profile_users p ON p.id_user = u.id WHERE p.id_profile = ?"
	queryGetProfileUsersAdds  = "SELECT id, name, role, status, (select id_profile from profile_users where id_user = u.id) FROM USERS u where u.id not in(select id_user from profile_users where id_profile = ?)"
	queryTotalProfileUsers    = "SELECT COUNT(*) as TOTAL FROM users u JOIN profile_users p ON p.id_user = u.id WHERE p.id_profile = ?"
	queryGetProfileAttendants = "SELECT id, name, role, status, (select id_profile from profile_users where id_user = u.id) FROM USERS u where u.id not in(select id_user from profile_users where id_profile = ?)"
	queryInsertProfileUser    = "INSERT INTO profile_users (id_profile, id_user) VALUES (?, ?)"
	queryUpdateProfileUser    = "UPDATE profile_users SET id_profile = ? WHERE id = ?"
	queryDeleteProfileUser    = "DELETE from profile_users WHERE id = ?"
	queryDeleteProfileMenu    = "DELETE from profile_menus WHERE id = ?"
	queryGetProfileUser       = "SELECT id, id_user, id_profile FROM profile_users WHERE id = ?"
	queryGetProfileMenu       = "SELECT id, id_menu, id_profile FROM profile_menus WHERE id = ?"
	queryGetProfileMenuFather = "SELECT id, id_menu, id_profile FROM profile_menus WHERE id_menu = ? and id_profile = ?"
	queryInsertProfileMenu    = "INSERT INTO profile_menus (id_menu, id_profile) VALUES (?, ?)"
	queryTotalProfileMenu     = "SELECT count(*) AS total FROM profile_menus WHERE id_menu = ? AND id_profile = ?"
	queryGetProfileRelation   = "SELECT m.id, m.name AS description, m.icon, m.link, m.parent, m.level FROM profile_users pu JOIN profile_menus pm ON pu.id_profile = pm.id_profile JOIN menus m ON pm.id_menu = m.id WHERE id_user = ? ORDER BY m.parent, m.menu_order"
)

func (p *Profile) GetProfileRelation() ([]BuildMenu, *rest_errors.RestErr) {
	stmt, err := stars_mysql.Client.Prepare(queryGetProfileRelation)

	if err != nil {
		logger.Error("error when trying to prepare get menus statement", err)
		return nil, rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, getErr := stmt.Query(p.Id)
	if getErr != nil {
		logger.Error("error when trying to get attendances", getErr)
		return nil, rest_errors.NewInternalServerError("database error")
	}
	defer rows.Close()

	results := make([]BuildMenu, 0)
	for rows.Next() {
		var bm BuildMenu

		if err := rows.Scan(&bm.Id, &bm.Description, &bm.Icon, &bm.Link, &bm.Parent, &bm.Level); err != nil {
			return nil, mysql_utils.ParseError(err)
		}
		results = append(results, bm)
	}

	return results, nil
}

func (pm *ProfileMenu) GetTotalProfileMenu(father int64) (*int64, *rest_errors.RestErr) {
	stmt, err := stars_mysql.Client.Prepare(queryTotalProfileMenu)

	if err != nil {
		logger.Error("error when trying to prepare total father statement", err)
		return nil, rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	totalRows := stmt.QueryRow(father, pm.IdProfile)
	var total int64

	if errTotalRows := totalRows.Scan(&total); errTotalRows != nil {
		logger.Error("error when trying to get total father", errTotalRows)
		return nil, rest_errors.NewInternalServerError("database error")
	}

	return &total, nil
}

func (p *Profile) GetProfile() *rest_errors.RestErr {
	stmt, err := stars_mysql.Client.Prepare(queryGetProfile)

	if err != nil {
		logger.Error("error when trying to prepare get bank statement", err)
		return rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(p.Id)

	if getErr := result.Scan(&p.Id, &p.Name, &p.ProfileCode); getErr != nil {
		logger.Error("error when trying to get bank", getErr)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func (p *ProfileUser) GetProfileUser() *rest_errors.RestErr {
	stmt, err := stars_mysql.Client.Prepare(queryGetProfileUser)

	if err != nil {
		logger.Error("error when trying to prepare get bank statement", err)
		return rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(p.IdUser)

	if getErr := result.Scan(&p.Id, &p.IdUser, &p.IdProfile); getErr != nil {
		logger.Error("error when trying to get bank (profile_user)", getErr)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func (p *ProfileMenu) GetProfileMenuFather() ([]ProfileMenu, *rest_errors.RestErr) {

	stmt, err := stars_mysql.Client.Prepare(queryGetProfileMenuFather)

	if err != nil {
		logger.Error("error when trying to prepare get menus statement", err)
		return nil, rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, getErr := stmt.Query(p.IdMenu, p.IdProfile)
	if getErr != nil {
		logger.Error("error when trying to get attendances", getErr)
		return nil, rest_errors.NewInternalServerError("database error")
	}
	defer rows.Close()

	results := make([]ProfileMenu, 0)
	for rows.Next() {
		var p ProfileMenu

		if err := rows.Scan(&p.Id, &p.IdMenu, &p.IdProfile); err != nil {
			return nil, mysql_utils.ParseError(err)
		}
		results = append(results, p)
	}

	return results, nil
}

func (p *ProfileMenu) GetProfileMenu() *rest_errors.RestErr {
	stmt, err := stars_mysql.Client.Prepare(queryGetProfileMenu)

	if err != nil {
		logger.Error("error when trying to prepare get profile_menu statement", err)
		return rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(p.IdMenu)

	if getErr := result.Scan(&p.Id, &p.IdMenu, &p.IdProfile); getErr != nil {
		logger.Error("error when trying to get profile_menu", getErr)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func buildQuery(query *string, queryTotal *string, filter *Filter) {

	concatQuery := ""

	if filter.Name != "" {
		concatQuery = concatQuery + " AND p.name LIKE '" + filter.Name + "%'"
	}
	if filter.ProfileCode != "" {
		concatQuery = concatQuery + " AND p.profile_code = '" + filter.ProfileCode + "%'"
	}

	if concatQuery != "" {
		*query = *query + concatQuery
		*queryTotal = *queryTotal + concatQuery
	}

	*query = *query + " LIMIT ?, ?"
}

func (p *Profile) GetProfiles(page int, itemsPerPage int, filter *Filter, userId int64) ([]Profile, *int, *rest_errors.RestErr) {
	query := queryGetProfiles
	queryTotal := queryTotalProfiles
	buildQuery(&query, &queryTotal, filter)

	stmt, err := stars_mysql.Client.Prepare(query)

	initialResult := (page - 1) * itemsPerPage

	if err != nil {
		logger.Error("error when trying to prepare get profiles statement", err)
		return nil, nil, rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	// aqio
	rows, getErr := stmt.Query(userId, initialResult, itemsPerPage)
	if getErr != nil {
		logger.Error("error when trying to get profiles", getErr)
		return nil, nil, rest_errors.NewInternalServerError("database error")
	}
	defer rows.Close()

	stmtTotalRows, err := stars_mysql.Client.Prepare(queryTotal)

	if err != nil {
		logger.Error("error when trying to prepare get total profiles rows statement", err)
		return nil, nil, rest_errors.NewInternalServerError("database error")
	}
	defer stmtTotalRows.Close()

	totalRows := stmtTotalRows.QueryRow(userId)
	var total int

	if errTotalRows := totalRows.Scan(&total); errTotalRows != nil {
		logger.Error("error when trying to get total profiles", errTotalRows)
		return nil, nil, rest_errors.NewInternalServerError("database error")
	}

	results := make([]Profile, 0)
	for rows.Next() {
		var pa Profile

		if err := rows.Scan(&pa.Id, &pa.Name, &pa.ProfileCode); err != nil {
			return nil, nil, mysql_utils.ParseError(err)
		}
		results = append(results, pa)
	}

	return results, &total, nil
}

func (p *Profile) Save() *rest_errors.RestErr {

	stmt, err := stars_mysql.Client.Prepare(queryInsertProfile)

	if err != nil {
		logger.Error("error when trying to prepare save instance statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()
	insertResult, saveErr := stmt.Exec(p.Name, p.ProfileCode)

	if saveErr != nil {
		logger.Error("error when trying to save profile", saveErr)
		return rest_errors.NewInternalServerError("database error")
	}

	profileId, err := insertResult.LastInsertId()

	if err != nil {
		logger.Error("error when trying to get last insert id after creating a new profile", err)
		return rest_errors.NewInternalServerError("database error")
	}

	p.Id = profileId

	return nil
}

func (pu *ProfileUser) SaveProfileUser() *rest_errors.RestErr {

	stmt, err := stars_mysql.Client.Prepare(queryInsertProfileUser)

	if err != nil {
		logger.Error("error when trying to prepare save pu instance statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()
	insertResult, saveErr := stmt.Exec(pu.IdProfile, pu.IdUser)

	if saveErr != nil {
		logger.Error("error when trying to save profile_user", saveErr)
		return rest_errors.NewInternalServerError("database error")
	}

	profileUserId, err := insertResult.LastInsertId()

	if err != nil {
		logger.Error("error when trying to get last insert id after creating a new profile_user", err)
		return rest_errors.NewInternalServerError("database error")
	}

	pu.Id = profileUserId

	return nil
}

func (pm *ProfileMenu) SaveProfileMenu() *rest_errors.RestErr {

	stmt, err := stars_mysql.Client.Prepare(queryInsertProfileMenu)

	if err != nil {
		logger.Error("error when trying to prepare save pu instance statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()
	insertResult, saveErr := stmt.Exec(pm.IdMenu, pm.IdProfile)

	if saveErr != nil {
		logger.Error("error when trying to save profile_menu", saveErr)
		return rest_errors.NewInternalServerError("database error")
	}

	profileMenuId, err := insertResult.LastInsertId()

	if err != nil {
		logger.Error("error when trying to get last insert id after creating a new profile_menu", err)
		return rest_errors.NewInternalServerError("database error")
	}

	pm.Id = profileMenuId

	return nil
}

func (pu *ProfileUser) UpdateProfileUser() *rest_errors.RestErr {

	stmt, err := stars_mysql.Client.Prepare(queryUpdateProfileUser)

	if err != nil {
		logger.Error("error when trying to prepare save pu instance statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	_, updateErr := stmt.Exec(&pu.IdProfile, &pu.Id)

	if updateErr != nil {
		logger.Error("error when trying to update profile_user", updateErr)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func (p *Profile) Update() *rest_errors.RestErr {

	stmt, err := stars_mysql.Client.Prepare(queryUpdateProfile)

	if err != nil {
		logger.Error("error when trying to prepare update profile statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	_, updateErr := stmt.Exec(&p.Name, &p.ProfileCode, &p.Id)

	if updateErr != nil {
		logger.Error("error when trying to update profile", updateErr)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func (p *Profile) Delete() *rest_errors.RestErr {

	stmt, err := stars_mysql.Client.Prepare(queryDeleteProfile)

	if err != nil {
		logger.Error("error when trying to prepare delete profile statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	_, deleteErr := stmt.Exec(p.Id)

	if deleteErr != nil {
		logger.Error("error when trying to delete profile", err)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func (pu *ProfileUser) DeleteProfileUser() *rest_errors.RestErr {

	stmt, err := stars_mysql.Client.Prepare(queryDeleteProfileUser)

	if err != nil {
		logger.Error("error when trying to prepare delete profile statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	_, deleteErr := stmt.Exec(pu.Id)

	if deleteErr != nil {
		logger.Error("error when trying to delete profile", err)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func (pu *ProfileMenu) DeleteProfileMenu() *rest_errors.RestErr {

	stmt, err := stars_mysql.Client.Prepare(queryDeleteProfileMenu)

	if err != nil {
		logger.Error("error when trying to prepare delete profile_menu statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	_, deleteErr := stmt.Exec(pu.Id)

	if deleteErr != nil {
		logger.Error("error when trying to delete profile_menu", err)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func (p *Profile) GetProfileUsers(page int, itemsPerPage int, filter *Filter, profileId int64) ([]User, *int, *rest_errors.RestErr) {

	query := queryGetProfileUsers
	queryTotal := queryTotalProfileUsers
	buildQuery(&query, &queryTotal, filter)

	stmt, err := stars_mysql.Client.Prepare(query)

	initialResult := (page - 1) * itemsPerPage

	if err != nil {
		logger.Error("error when trying to prepare get users statement", err)
		return nil, nil, rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, getErr := stmt.Query(profileId, initialResult, itemsPerPage)
	if getErr != nil {
		logger.Error("error when trying to get users", getErr)
		return nil, nil, rest_errors.NewInternalServerError("database error")
	}
	defer rows.Close()

	stmtTotalRows, err := stars_mysql.Client.Prepare(queryTotal)

	if err != nil {
		logger.Error("error when trying to prepare get total users rows statement", err)
		return nil, nil, rest_errors.NewInternalServerError("database error")
	}
	defer stmtTotalRows.Close()

	totalRows := stmtTotalRows.QueryRow(profileId)
	var total int

	if errTotalRows := totalRows.Scan(&total); errTotalRows != nil {
		logger.Error("error when trying to get total profiles_users", errTotalRows)
		return nil, nil, rest_errors.NewInternalServerError("database error")
	}

	if getErr != nil {
		logger.Error("error when trying to get profiles_users", getErr)
		return nil, nil, rest_errors.NewInternalServerError("database error")
	}

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Name, &user.Role, &user.Status, &user.IdProfile); err != nil {
			return nil, nil, mysql_utils.ParseError(err)
		}

		results = append(results, user)
	}

	return results, &total, nil
}

func (p *Profile) GetProfileUsersAdds(page int, itemsPerPage int, filter *Filter, profileId int64) ([]User, *rest_errors.RestErr) {

	query := queryGetProfileUsersAdds
	queryTotal := queryTotalProfileUsers
	buildQuery(&query, &queryTotal, filter)

	stmt, err := stars_mysql.Client.Prepare(query)

	initialResult := (page - 1) * itemsPerPage

	if err != nil {
		logger.Error("error when trying to prepare get users statement", err)
		return nil, rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, getErr := stmt.Query(profileId, initialResult, itemsPerPage)
	if getErr != nil {
		logger.Error("error when trying to get profiles_users", getErr)
		return nil, rest_errors.NewInternalServerError("database error")
	}
	defer rows.Close()

	if getErr != nil {
		logger.Error("error when trying to get profiles_users", getErr)
		return nil, rest_errors.NewInternalServerError("database error")
	}

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Name, &user.Role, &user.Status, &user.IdProfile); err != nil {
			return nil, mysql_utils.ParseError(err)
		}

		results = append(results, user)
	}

	return results, nil
}

func (user *User) GetProfileAttendants(search string, profileId int64) ([]User, *rest_errors.RestErr) {

	query := queryGetProfileAttendants + " AND name LIKE '%" + search + "%'"

	stmt, err := stars_mysql.Client.Prepare(query)

	if err != nil {
		logger.Error("error when trying to prepare get attendances statement", err)
		return nil, rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, getErr := stmt.Query(profileId)
	if getErr != nil {
		logger.Error("error when trying to get attendances", getErr)
		return nil, rest_errors.NewInternalServerError("database error")
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Name, &user.Role, &user.Status, &user.IdProfile); err != nil {
			return nil, mysql_utils.ParseError(err)
		}
		results = append(results, user)
	}

	return results, nil
}
