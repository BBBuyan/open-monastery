package db

import (
	"errors"
	"fmt"

	"example.com/open-monastery/internal/app/apperr"
	mysqlDriver "github.com/go-sql-driver/mysql"
)

func translateWriteError(operation string, err error) error {
	if mysqlErr, ok := errors.AsType[*mysqlDriver.MySQLError](err); ok {
		switch mysqlErr.Number {
		case 1062:
			return apperr.NewConflict(
				"resource already exists",
				err,
			)
		case 1452:
			return apperr.NewNotFound(
				"referenced resource not found",
				err,
			)
		}
	}

	return fmt.Errorf("%s: %w", operation, err)
}
