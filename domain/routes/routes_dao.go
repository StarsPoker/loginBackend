package routes

import (
	"github.com/StarsPoker/loginBackend/logger"

	"github.com/StarsPoker/loginBackend/datasources/mysql/stars_mysql"
	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
	"github.com/StarsPoker/loginBackend/utils/mysql_utils"
)

const (
	queryGetRoute    = "SELECT id, name, type, menu_id FROM routes WHERE id = ?"
	queryTotalRoutes = "SELECT COUNT(*) as TOTAL FROM routes r WHERE 1 = 1"
	queryGetRoutes   = "SELECT r.id, r.name, r.type, r.menu_id, m.name AS menu_string FROM routes r LEFT JOIN menus m ON m.id = r.menu_id WHERE 1 = 1"
	queryInsertRoute = "INSERT INTO routes (name, type, menu_id) VALUES (?, ?, ?)"
	queryUpdateRoute = "UPDATE routes SET name = ?, type = ?, menu_id = ? WHERE id = ?"
	queryDeleteRoute = "DELETE from routes WHERE id = ?"
)

func (r *Route) GetRoute() *rest_errors.RestErr {
	stmt, err := stars_mysql.Client.Prepare(queryGetRoute)

	if err != nil {
		logger.Error("error when trying to prepare get bank statement", err)
		return rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(r.Id)

	if getErr := result.Scan(&r.Id, &r.Name, &r.Type, &r.MenuId); getErr != nil {
		logger.Error("error when trying to get bank", getErr)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func (r *Route) GetRoutes(page int, itemsPerPage int, filter *Filter) ([]Route, *int, *rest_errors.RestErr) {
	query := queryGetRoutes
	queryTotal := queryTotalRoutes
	buildQuery(&query, &queryTotal, filter)

	stmt, err := stars_mysql.Client.Prepare(query)

	initialResult := (page - 1) * itemsPerPage

	if err != nil {
		logger.Error("error when trying to prepare get routes statement", err)
		return nil, nil, rest_errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, getErr := stmt.Query(initialResult, itemsPerPage)
	if getErr != nil {
		logger.Error("error when trying to get routes", getErr)
		return nil, nil, rest_errors.NewInternalServerError("database error")
	}
	defer rows.Close()

	stmtTotalRows, err := stars_mysql.Client.Prepare(queryTotal)

	if err != nil {
		logger.Error("error when trying to prepare get total routes rows statement", err)
		return nil, nil, rest_errors.NewInternalServerError("database error")
	}
	defer stmtTotalRows.Close()

	totalRows := stmtTotalRows.QueryRow()
	var total int

	if errTotalRows := totalRows.Scan(&total); errTotalRows != nil {
		logger.Error("error when trying to get total routes", errTotalRows)
		return nil, nil, rest_errors.NewInternalServerError("database error")
	}

	results := make([]Route, 0)
	for rows.Next() {
		var ro Route

		if err := rows.Scan(&ro.Id, &ro.Name, &ro.Type, &ro.MenuId, &ro.MenuSt); err != nil {
			return nil, nil, mysql_utils.ParseError(err)
		}
		results = append(results, ro)
	}

	return results, &total, nil
}

func buildQuery(query *string, queryTotal *string, filter *Filter) {

	concatQuery := ""

	if filter.Name != "" {
		concatQuery = concatQuery + " AND r.name LIKE '" + filter.Name + "%'"
	}
	if filter.Type != "" {
		concatQuery = concatQuery + " AND r.type = '" + filter.Type + "'"
	}
	if filter.MenuId != "" {
		concatQuery = concatQuery + " AND r.menu_id = '" + filter.MenuId + "'"
	}

	if concatQuery != "" {
		*query = *query + concatQuery
		*queryTotal = *queryTotal + concatQuery
	}

	*query = *query + " LIMIT ?, ?"
}

func (r *Route) Save() *rest_errors.RestErr {

	stmt, err := stars_mysql.Client.Prepare(queryInsertRoute)

	if err != nil {
		logger.Error("error when trying to prepare save instance statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()
	insertResult, saveErr := stmt.Exec(r.Name, r.Type, r.MenuId)

	if saveErr != nil {
		logger.Error("error when trying to save route", saveErr)
		return rest_errors.NewInternalServerError("database error")
	}

	routeId, err := insertResult.LastInsertId()

	if err != nil {
		logger.Error("error when trying to get last insert id after creating a new profile", err)
		return rest_errors.NewInternalServerError("database error")
	}

	r.Id = routeId

	return nil
}

func (r *Route) Update() *rest_errors.RestErr {

	stmt, err := stars_mysql.Client.Prepare(queryUpdateRoute)

	if err != nil {
		logger.Error("error when trying to prepare update route statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	_, updateErr := stmt.Exec(&r.Name, &r.Type, &r.MenuId, &r.Id)

	if updateErr != nil {
		logger.Error("error when trying to update route", updateErr)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}

func (r *Route) Delete() *rest_errors.RestErr {

	stmt, err := stars_mysql.Client.Prepare(queryDeleteRoute)

	if err != nil {
		logger.Error("error when trying to prepare delete route statement", err)
		return rest_errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	_, deleteErr := stmt.Exec(r.Id)

	if deleteErr != nil {
		logger.Error("error when trying to delete route", err)
		return rest_errors.NewInternalServerError("database error")
	}

	return nil
}
