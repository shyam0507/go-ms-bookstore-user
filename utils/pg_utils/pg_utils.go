package pgutils

import (
	"go-ms-bookstore-user/utils/errors"
	"strings"

	"github.com/lib/pq"
)

const (
	noRows = "sql: no rows in result set"
)

func ParseError(err error) *errors.RestErr {
	pgErr, ok := err.(*pq.Error)

	if !ok {
		if strings.Contains(err.Error(), noRows) {
			return errors.NewNotFoundError("no record matching given id!")
		}

		return errors.NewInternalServerError("Error parsing database response")
	}

	switch pgErr.Code {
	case "23505":
		return errors.NewBadRequestError("Invalid data")
	}
	return errors.NewInternalServerError("Error processing request")

}
