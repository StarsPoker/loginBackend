package users

import (
	"fmt"
	"strings"

	"github.com/StarsPoker/loginBackend/logger"

	"github.com/StarsPoker/loginBackend/utils/mysql_utils"

	"github.com/StarsPoker/loginBackend/datasources/mysql/stars_mysql"
	"github.com/StarsPoker/loginBackend/utils/date_utils"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
)

const (
	errorNoRows                 = "no rows in result set"
	queryGetExternalAccess      = "SELECT u.id, u.external_access FROM users u WHERE u.id = ?"
	queryGetUser                = "SELECT u.id, u.name, u.email, u.password, p.profile_code as role, u.status, DATE_FORMAT(date_created, '%d/%m/%Y %k:%i'), u.instance_id, u.default_password, p.name, u.authenticator_configured, u.otp_secret FROM users u LEFT JOIN profile_users pu ON pu.id_user = u.id LEFT JOIN profiles p ON p.id = pu.id_profile WHERE u.id = ?"
	queryTotalUsers             = "SELECT COUNT(*) as TOTAL FROM users u WHERE 1 = 1"
	queryGetUsers               = "SELECT u.id, u.name, u.email, u.contact, u.password, u.status, DATE_FORMAT(date_created, '%d/%m/%Y %k:%i') date_created, u.instance_id, u.default_password, i.name as instance_name, u.inscription, u.authenticator_configured FROM users u LEFT JOIN instances i ON u.instance_id = i.id WHERE 1 = 1"
	queryGetAttendants          = "SELECT id, name,  role, status FROM users WHERE 1 = 1"
	queryFindByEmailAndPassword = "SELECT id, name, email, role, contact, status, DATE_FORMAT(date_created, '%d/%m/%Y %k:%i') date_created, password, inscription from users  WHERE email = ? AND status = ?"
	queryInsertUser             = "INSERT INTO users (name, email, contact, password, status, date_created, instance_id, default_password, inscription, otp_secret, authenticator_configured) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	queryUpdateUser             = "UPDATE users SET email = ?, status = ?, instance_id = ?, name = ?, contact = ?, inscription = ?, otp_secret = ?, authenticator_configured = ? WHERE id = ?"
	queryUpdateUserName         = "UPDATE users SET name = ? WHERE id = ?"
	queryUpdateUserEmail        = "UPDATE users SET email = ? WHERE id = ?"
	queryChangePassword         = "UPDATE users SET password = ?, default_password = 0 WHERE id = ?"
	queryDeleteUser             = "DELETE FROM users WHERE id = ?"
)

var (
	usersDB = make(map[int64]*User)
)

func buildQuery(query *string, queryTotal *string, filter *Filter) {

	concatQuery := ""

	if filter.Role != "" {
		concatQuery = concatQuery + " AND u.role = " + filter.Role
	}

	if filter.Name != "" {
		concatQuery = concatQuery + " AND u.name LIKE '" + filter.Name + "%'"
	}

	if filter.Email != "" {
		concatQuery = concatQuery + " AND u.email LIKE '" + filter.Email + "%'"
	}

	if filter.Club != "" {
		concatQuery = concatQuery + " AND u.instance_id = " + filter.Club
	}

	if filter.Status != "" {
		concatQuery = concatQuery + " AND u.status = " + filter.Status
	}

	if filter.DefaultPassword != "" {
		concatQuery = concatQuery + " AND u.default_password = " + filter.DefaultPassword
	}

	if concatQuery != "" {
		*query = *query + concatQuery
		*queryTotal = *queryTotal + concatQuery
	}

	if filter.SortBy != "" {
		*query = *query + " ORDER BY u." + filter.SortBy
		if filter.SortDesc == "true" {
			*query = *query + " desc"
		}
	} else {
		*query = *query + " ORDER BY u.status, u.instance_id"
	}
	*query = *query + " LIMIT ?, ?"
}

func (user *User) ValidateExternalAccess(user_id int64) *rest_errors.RestErr {
	stmt, err := stars_mysql.Client.Prepare(queryGetExternalAccess)
	if err != nil {
		logger.Error("error when trying to prepare get external access statement", err)
		return rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user_id)
	if getErr := result.Scan(&user.Id, &user.ExternalAccess); getErr != nil {
		logger.Error("error when trying to get user external access", getErr)
		return rest_errors.NewInternalServerError("database error")
	}

	if user.ExternalAccess == 0 {
		logger.Error("user with external access blocked", nil)
		return rest_errors.NewInternalServerError("database error (external access blocked)")
	}

	return nil
}

func (user *User) FindByEmailAndPassword() *rest_errors.RestErr {
	stmt, err := stars_mysql.Client.Prepare(queryFindByEmailAndPassword)

	if err != nil {
		logger.Error("error when trying to get user by email and password", err)
		return rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Email, user.Status)

	if getErr := result.Scan(&user.Id, &user.Name, &user.Email, &user.Role, &user.Contact, &user.Status, &user.DateCreated, &user.Password, &user.Inscription); getErr != nil {
		if strings.Contains(getErr.Error(), mysql_utils.ErrorNoRows) {
			return rest_errors.NewNotFoundError("invalid user credentials")
		}
		logger.Error("error when trying to get user login", getErr)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) GetUsers(page int, itemsPerPage int, filter *Filter) ([]User, *int, *rest_errors.RestErr) {

	query := queryGetUsers
	queryTotal := queryTotalUsers
	buildQuery(&query, &queryTotal, filter)
	fmt.Println(query)
	stmt, err := stars_mysql.Client.Prepare(query)

	initialResult := (page - 1) * itemsPerPage

	if err != nil {
		logger.Error("error when trying to prepare get users statement", err)
		return nil, nil, rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, getErr := stmt.Query(initialResult, itemsPerPage)
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

	totalRows := stmtTotalRows.QueryRow()
	var total int

	if errTotalRows := totalRows.Scan(&total); errTotalRows != nil {
		logger.Error("error when trying to get total users", errTotalRows)
		return nil, nil, rest_errors.NewInternalServerError("database error")
	}

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Contact, &user.Password, &user.Status, &user.DateCreated,
			&user.InstanceId, &user.DefaultPassword, &user.InstanceName, &user.Inscription, &user.AuthenticatorConfigured); err != nil {
			fmt.Println(err)
			return nil, nil, mysql_utils.ParseError(err)
		}

		results = append(results, user)
	}

	return results, &total, nil
}

func (user *User) GetAttendants(search string) ([]User, *rest_errors.RestErr) {

	query := queryGetAttendants + " AND name LIKE '%" + search + "%'"

	stmt, err := stars_mysql.Client.Prepare(query)

	if err != nil {
		logger.Error("error when trying to prepare get attendances statement", err)
		return nil, rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, getErr := stmt.Query()
	if getErr != nil {
		logger.Error("error when trying to get attendances", getErr)
		return nil, rest_errors.NewInternalServerError("database error")
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Name, &user.Role, &user.Status); err != nil {
			return nil, mysql_utils.ParseError(err)
		}
		results = append(results, user)
	}

	return results, nil
}

func (user *User) GetUser() *rest_errors.RestErr {
	stmt, err := stars_mysql.Client.Prepare(queryGetUser)

	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)

	if getErr := result.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Role, &user.Status, &user.DateCreated,
		&user.InstanceId, &user.DefaultPassword, &user.ProfileAccess, &user.AuthenticatorConfigured, &user.OTPSecret); getErr != nil {
		logger.Error("error when trying to get user", getErr)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) Save() *rest_errors.RestErr {
	user.DateCreated = date_utils.GetNowDBFormat()

	stmt, err := stars_mysql.Client.Prepare(queryInsertUser)

	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(user.Name, user.Email, user.Contact, user.Password, user.Status, user.DateCreated,
		user.InstanceId, user.DefaultPassword, user.Inscription, user.OTPSecret, user.AuthenticatorConfigured)

	if saveErr != nil {
		logger.Error("error when trying to save user", saveErr)
		return rest_errors.NewInternalServerError("database error")
	}

	userId, err := insertResult.LastInsertId()

	if err != nil {
		logger.Error("error when trying to get last insert id after creating a new user", err)
		return rest_errors.NewInternalServerError("database error")
	}

	user.Id = userId

	return nil
}

func (user *User) Update() *rest_errors.RestErr {

	stmt, err := stars_mysql.Client.Prepare(queryUpdateUser)

	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	_, updateErr := stmt.Exec(user.Email, user.Status, user.InstanceId, user.Name, user.Contact, user.Inscription, user.OTPSecret, user.AuthenticatorConfigured, user.Id)

	if updateErr != nil {
		logger.Error("error when trying to update user", updateErr)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) UpdateUserName() *rest_errors.RestErr {

	stmt, err := stars_mysql.Client.Prepare(queryUpdateUserName)

	if err != nil {
		logger.Error("error when trying to prepare update username statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	_, updateErr := stmt.Exec(user.Name, user.Id)

	if updateErr != nil {
		logger.Error("error when trying to update username", updateErr)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) UpdateUserEmail() *rest_errors.RestErr {

	stmt, err := stars_mysql.Client.Prepare(queryUpdateUserEmail)

	if err != nil {
		logger.Error("error when trying to prepare update useremail statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	_, updateErr := stmt.Exec(user.Email, user.Id)

	if updateErr != nil {
		logger.Error("error when trying to update useremail", updateErr)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) ChangePassword() *rest_errors.RestErr {

	// Se o usuário já fez alguma alteração na senha, a query desta função mudará o valor de default_password para 0. Se não, o valor permanece como 1.

	stmt, err := stars_mysql.Client.Prepare(queryChangePassword)

	if err != nil {
		logger.Error("error when trying to prepare change password user statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	_, updateErr := stmt.Exec(user.Password, user.Id)

	if updateErr != nil {
		logger.Error("error when trying to change password user", updateErr)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) Delete() *rest_errors.RestErr {

	stmt, err := stars_mysql.Client.Prepare(queryDeleteUser)

	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	_, deleteErr := stmt.Exec(user.Id)

	if deleteErr != nil {
		logger.Error("error when trying to delete user", err)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}
