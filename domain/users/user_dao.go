package users

import (
	"strings"

	"github.com/StarsPoker/loginBackend/logger"

	"github.com/StarsPoker/loginBackend/utils/mysql_utils"

	"github.com/StarsPoker/loginBackend/datasources/mysql/stars_mysql"
	"github.com/StarsPoker/loginBackend/utils/date_utils"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
)

const (
	errorNoRows                 = "no rows in result set"
	queryGetUser                = "SELECT id, name, email, password, role, status, DATE_FORMAT(date_created, '%d/%m/%Y %k:%i') FROM users WHERE id = ?"
	queryGetUsers               = "SELECT id, name, email, password, role, status, DATE_FORMAT(date_created, '%d/%m/%Y %k:%i') date_created FROM users WHERE 1 = 1"
	queryFindByEmailAndPassword = "SELECT id, name, email, role, status, DATE_FORMAT(date_created, '%d/%m/%Y %k:%i') date_created from users WHERE email = ? AND password = ? AND status = ?"
	queryInsertUser             = "INSERT INTO USERS (name, email, password, role, status, date_created) VALUES (?, ?, ?, ?, ?, ?)"
	queryUpdateUser             = "UPDATE USERS SET email = ?, status = ?, role = ? WHERE id = ?"
	queryDeleteUser             = "DELETE FROM USERS WHERE id = ?"
)

var (
	usersDB = make(map[int64]*User)
)

func (user *User) FindByEmailAndPassword() *rest_errors.RestErr {
	stmt, err := stars_mysql.Client.Prepare(queryFindByEmailAndPassword)

	if err != nil {
		logger.Error("error when trying to get user by email and password", err)
		return rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Email, user.Password, user.Status)

	if getErr := result.Scan(&user.Id, &user.Name, &user.Email, &user.Role, &user.Status, &user.DateCreated); getErr != nil {
		if strings.Contains(getErr.Error(), mysql_utils.ErrorNoRows) {
			return rest_errors.NewNotFoundError("invalid user credentials")
		}
		logger.Error("error when trying to get user login", getErr)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) GetUsers() ([]User, *rest_errors.RestErr) {
	stmt, err := stars_mysql.Client.Prepare(queryGetUsers)

	if err != nil {
		logger.Error("error when trying to prepare get users statement", err)
		return nil, rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, getErr := stmt.Query()

	if getErr != nil {
		logger.Error("error when trying to get users", getErr)
		return nil, rest_errors.NewInternalServerError("database error")
	}

	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Role, &user.Status, &user.DateCreated); err != nil {
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

	if getErr := result.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Role, &user.Status, &user.DateCreated); getErr != nil {
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

	insertResult, saveErr := stmt.Exec(user.Name, user.Email, user.Password, user.Role, user.Status, user.DateCreated)

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

	_, updateErr := stmt.Exec(user.Email, user.Status, user.Role, user.Id)

	if updateErr != nil {
		logger.Error("error when trying to update user", updateErr)
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
