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
	queryGetChildrens         = "SELECT m.id, m.name, m.icon, m.link, m.parent, m.level, m.menu_order, m.profile_father, coalesce(pm.id, 0) as id_relation, IF(pm.id IS NULL,0, 1) as has_relation FROM menus m LEFT JOIN profile_menus pm ON pm.id_menu = m.id AND pm.id_profile = ? WHERE parent = ? ORDER BY menu_order, m.name"
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
	queryGetChildrenSearch    = "SELECT m.id, m.name, m.parent, m.link FROM menus m WHERE 2 = 2"
	queryGetUserPermission    = "SELECT COUNT(*) FROM profile_users pu JOIN profile_menus pm ON pu.id_profile = pm.id_profile JOIN menus m ON pm.id_menu = m.id WHERE id_user = ? AND m.link = ? ORDER BY m.parent, m.menu_order"
	queryGetProfilesRelation  = "SELECT p.id, p.name, pm.id as id_relation, IF(m.name IS NULL,0, 1) as has_relation FROM profiles p LEFT JOIN profile_menus pm ON pm.id_profile = p.id AND pm.id_menu = ? LEFT JOIN menus m ON m.id = pm.id_menu  WHERE p.profile_code < 100"
)

func (pp *ProfilePermission) GetUserPermission() (*int, *rest_errors.RestErr) {
	stmt, err := stars_mysql.Client.Prepare(queryGetUserPermission)

	if err != nil {
		logger.Error("error when trying to prepare total father statement", err)
		return nil, rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	totalRows := stmt.QueryRow(pp.UserId, pp.MenuName)
	var total int

	if errTotalRows := totalRows.Scan(&total); errTotalRows != nil {
		logger.Error("error when trying to get total father", errTotalRows)
		return nil, rest_errors.NewInternalServerError("database error")
	}

	return &total, nil
}

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
		logger.Error("error when trying to get menu", getErr)
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

func (menu *Menus) GetChildrenSearch(search string) ([]Menu, *rest_errors.RestErr) {

	query := queryGetChildrenSearch + " AND m.name LIKE '%" + search + "%'"

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

	results := make([]Menu, 0)
	for rows.Next() {
		var menu Menu
		if err := rows.Scan(&menu.Id, &menu.Name, &menu.Parent, &menu.Link); err != nil {
			return nil, mysql_utils.ParseError(err)
		}
		results = append(results, menu)
	}

	return results, nil
}

func (menu *Menus) GetProfilesRelation(menuId int64) ([]ProfileRelation, *rest_errors.RestErr) {

	stmt, err := stars_mysql.Client.Prepare(queryGetProfilesRelation)

	if err != nil {
		logger.Error("error when trying to prepare get profilesRelation statement", err)
		return nil, rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, getErr := stmt.Query(menuId)
	if getErr != nil {
		logger.Error("error when trying to get profilesRelation", getErr)
		return nil, rest_errors.NewInternalServerError("database error")
	}
	defer rows.Close()

	results := make([]ProfileRelation, 0)
	for rows.Next() {
		var p ProfileRelation

		if err := rows.Scan(&p.MenuId, &p.MenuName, &p.IdRelation, &p.HasRelation); err != nil {
			return nil, mysql_utils.ParseError(err)
		}
		results = append(results, p)
	}

	return results, nil
}
