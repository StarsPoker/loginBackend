package menus

import (
	"github.com/StarsPoker/loginBackend/datasources/mysql/stars_mysql"
	"github.com/StarsPoker/loginBackend/logger"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
	"github.com/StarsPoker/loginBackend/utils/mysql_utils"
)

const (
	errorNoRows               = "no rows in result set"
	queryGetMenu              = "SELECT m.id, m.name, m.icon, m.link, m.parent, m.level, m.menu_order FROM menus m WHERE id = ?"
	queryGetMenus             = "SELECT DISTINCT m.id, m.name, m.icon, m.link, m.parent, m.level, m.menu_order, coalesce(pm.id, 0) as id_relation, IF(pm.id IS NULL,0, 1) as has_relation, IF(mm.id IS NULL, 0, 1) as has_chield FROM menus m LEFT JOIN menus mm  ON mm.parent = m.id LEFT JOIN profile_menus pm ON pm.id_menu = m.id AND pm.id_profile = ? WHERE m.level = 1 ORDER BY menu_order"
	queryGetChildrens         = "SELECT m.id, m.name, m.icon, m.link, m.parent, m.level, m.menu_order, m.profile_father, coalesce(pm.id, 0) as id_relation, IF(pm.id IS NULL,0, 1) as has_relation FROM menus m LEFT JOIN profile_menus pm ON pm.id_menu = m.id AND pm.id_profile = ? WHERE parent = ? ORDER BY menu_order"
	queryinsertMenu           = "INSERT INTO menus (name, icon, link, parent, level, menu_order) VALUES (?, ?, ?, ?, ?, ?)"
	queryUpdateMenu           = "UPDATE menus SET name = ?, icon = ?, link = ? WHERE id = ?"
	queryDeleteMenu           = "DELETE FROM menus WHERE id = ?"
	queryDeleteChildren       = "DELETE FROM menus WHERE parent = ?"
	queryUpdateOrder          = "UPDATE menus SET menu_order = (menu_order - 1) WHERE menu_order > ?"
	queryUpdateOrderUpFirst   = "UPDATE menus SET menu_order = (menu_order + 1) WHERE menu_order = (? - 1)"
	queryUpdateOrderUpNext    = "UPDATE menus SET menu_order = (menu_order - 1) WHERE id = ?"
	queryUpdateOrderDownFirst = "UPDATE menus SET menu_order = (menu_order - 1) WHERE menu_order = (? + 1)"
	queryUpdateOrderDownNext  = "UPDATE menus SET menu_order = (menu_order + 1) WHERE id = ?"
	queryGetMaxOrder          = "SELECT count(*) as total from menus m where m.level = 1"
)

func (me *Menu) GetMenus() ([]Menu, *rest_errors.RestErr) {

	stmt, err := stars_mysql.Client.Prepare(queryGetMenus)

	if err != nil {
		logger.Error("error when trying to prepare get menus statement", err)
		return nil, rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, getErr := stmt.Query(me.ProfileFather)
	if getErr != nil {
		logger.Error("error when trying to get attendances", getErr)
		return nil, rest_errors.NewInternalServerError("database error")
	}
	defer rows.Close()

	results := make([]Menu, 0)
	for rows.Next() {
		var m Menu

		if err := rows.Scan(&m.Id, &m.Name, &m.Icon, &m.Link, &m.Parent, &m.Level, &m.Order, &m.IdRelation, &m.HasRelation, &m.HasChield); err != nil {
			return nil, mysql_utils.ParseError(err)
		}
		results = append(results, m)
	}

	return results, nil
}

func (me *Menu) GetChildrens() ([]Menu, *rest_errors.RestErr) {

	stmt, err := stars_mysql.Client.Prepare(queryGetChildrens)

	if err != nil {
		logger.Error("error when trying to prepare get menus statement", err)
		return nil, rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, getErr := stmt.Query(me.ProfileFather, me.Id)
	if getErr != nil {
		logger.Error("error when trying to get attendances", getErr)
		return nil, rest_errors.NewInternalServerError("database error")
	}
	defer rows.Close()

	results := make([]Menu, 0)
	for rows.Next() {
		var m Menu

		if err := rows.Scan(&m.Id, &m.Name, &m.Icon, &m.Link, &m.Parent, &m.Level, &m.Order, &m.ProfileFather, &m.IdRelation, &m.HasRelation); err != nil {
			return nil, mysql_utils.ParseError(err)
		}
		results = append(results, m)
	}

	return results, nil
}

func (me *Menu) Save() *rest_errors.RestErr {

	stmt, err := stars_mysql.Client.Prepare(queryinsertMenu)

	if err != nil {
		logger.Error("error when trying to prepare save instance statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()
	insertResult, saveErr := stmt.Exec(me.Name, me.Icon, me.Link, me.Parent, me.Level, me.Order)

	if saveErr != nil {
		logger.Error("error when trying to save menu", saveErr)
		return rest_errors.NewInternalServerError("database error")
	}

	menuId, err := insertResult.LastInsertId()

	if err != nil {
		logger.Error("error when trying to get last insert id after creating a new menu", err)
		return rest_errors.NewInternalServerError("database error")
	}

	me.Id = menuId

	return nil
}

func (m *Menu) GetMenu() *rest_errors.RestErr {
	stmt, err := stars_mysql.Client.Prepare(queryGetMenu)

	if err != nil {
		logger.Error("error when trying to prepare get bank statement (menu)", err)
		return rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(m.Id)

	if getErr := result.Scan(&m.Id, &m.Name, &m.Icon, &m.Link, &m.Parent, &m.Level, &m.Order); getErr != nil {
		logger.Error("error when trying to get bank (menu)", getErr)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func (m *Menu) GetMaxOrder() *int64 {
	stmt, err := stars_mysql.Client.Prepare(queryGetMaxOrder)

	if err != nil {
		logger.Error("error when trying to prepare get statement (menu)", err)
		return nil
	}
	defer stmt.Close()

	result := stmt.QueryRow()

	var maxOrder *int64

	if getErr := result.Scan(&maxOrder); getErr != nil {
		logger.Error("error when trying to get (Max Order)", getErr)
		return nil
	}

	return maxOrder
}

func (m *Menu) Delete() *rest_errors.RestErr {

	stmt, err := stars_mysql.Client.Prepare(queryDeleteMenu)

	if err != nil {
		logger.Error("error when trying to prepare delete profile statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	_, deleteErr := stmt.Exec(m.Id)

	if deleteErr != nil {
		logger.Error("error when trying to delete profile", err)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func (m *Menu) DeleteChildren() *rest_errors.RestErr {

	stmt, err := stars_mysql.Client.Prepare(queryDeleteChildren)

	if err != nil {
		logger.Error("error when trying to prepare delete profile statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	_, deleteErr := stmt.Exec(m.Id)

	if deleteErr != nil {
		logger.Error("error when trying to delete profile", err)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func (m *Menu) Update() *rest_errors.RestErr {

	stmt, err := stars_mysql.Client.Prepare(queryUpdateMenu)

	if err != nil {
		logger.Error("error when trying to prepare update menu statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	_, updateErr := stmt.Exec(&m.Name, &m.Icon, &m.Link, &m.Id)

	if updateErr != nil {
		logger.Error("error when trying to update menu", updateErr)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func (m *Menu) UpdateOrder() *rest_errors.RestErr {

	stmt, err := stars_mysql.Client.Prepare(queryUpdateOrder)

	if err != nil {
		logger.Error("error when trying to prepare update menu_order statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	_, updateErr := stmt.Exec(m.Order)

	if updateErr != nil {
		logger.Error("error when trying to update menu_oder", updateErr)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func (m *Menu) UpdateOrderUpFirst() *rest_errors.RestErr {

	stmt, err := stars_mysql.Client.Prepare(queryUpdateOrderUpFirst)

	if err != nil {
		logger.Error("error when trying to prepare update menu_order statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	_, updateErr := stmt.Exec(m.Order)

	if updateErr != nil {
		logger.Error("error when trying to update menu_oder", updateErr)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func (m *Menu) UpdateOrderUpNext() *rest_errors.RestErr {

	stmt, err := stars_mysql.Client.Prepare(queryUpdateOrderUpNext)

	if err != nil {
		logger.Error("error when trying to prepare update menu_order statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	_, updateErr := stmt.Exec(m.Id)

	if updateErr != nil {
		logger.Error("error when trying to update menu_oder", updateErr)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func (m *Menu) UpdateOrderDownFirst() *rest_errors.RestErr {

	stmt, err := stars_mysql.Client.Prepare(queryUpdateOrderDownFirst)

	if err != nil {
		logger.Error("error when trying to prepare update menu_order statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	_, updateErr := stmt.Exec(m.Order)

	if updateErr != nil {
		logger.Error("error when trying to update menu_oder", updateErr)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func (m *Menu) UpdateOrderDownNext() *rest_errors.RestErr {

	stmt, err := stars_mysql.Client.Prepare(queryUpdateOrderDownNext)

	if err != nil {
		logger.Error("error when trying to prepare update menu_order statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	_, updateErr := stmt.Exec(m.Id)

	if updateErr != nil {
		logger.Error("error when trying to update menu_oder", updateErr)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}
