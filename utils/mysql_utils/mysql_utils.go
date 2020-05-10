package mysql_utils

import (
	"fmt"
	"strings"

	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
	"github.com/go-sql-driver/mysql"
)

const (
	errorNoRows = "no rows in result set"
)

func ParseError(err error) *rest_errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), errorNoRows) {
			return rest_errors.NewNotFoundError("no record matching given id")
		}
		fmt.Println(err.Error())
		return rest_errors.NewInternalServerError("error parsing database response")
	}

	switch sqlErr.Number {
	case 1862:
		return rest_errors.NewBadRequestError("invalid data")
	}
	fmt.Println(sqlErr.Number)
	return rest_errors.NewInternalServerError("error processing request")
}
