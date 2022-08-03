package external_access_dao

import (
	"github.com/StarsPoker/loginBackend/datasources/mysql/stars_mysql"
	"github.com/StarsPoker/loginBackend/logger"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
)

const (
	queryGetMenu = "SELECT m.id, m.name, m.icon, m.link, m.parent, m.level, m.menu_order FROM menus m WHERE id = ?"
)

func (m *Menu) GetMenu() *rest_errors.RestErr {
	stmt, err := stars_mysql.Client.Prepare(queryGetMenu)

	if err != nil {
		logger.Error("error when trying to prepare get bank statement (menu)", err)
		return rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(m.Id)

	if getErr := result.Scan(&m.Id, &m.Name, &m.Icon, &m.Link, &m.Parent, &m.Level, &m.Order); getErr != nil {
		logger.Error("error when trying to get menu", getErr)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}
