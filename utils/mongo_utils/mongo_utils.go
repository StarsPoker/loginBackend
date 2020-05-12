package mongo_utils

import (
	"strings"

	"github.com/StarsPoker/loginBackend/utils/errors/rest_errors"
)

const (
	errorNoRows   = "no documents in result"
	errorNotFound = "not found"
)

func ParseError(err error) *rest_errors.RestErr {

	if strings.Contains(err.Error(), errorNoRows) {
		return rest_errors.NewNotFoundError("no record matching given id")
	}

	return rest_errors.NewInternalServerError("error parsing database response")
}
