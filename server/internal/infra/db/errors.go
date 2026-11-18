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
				"already exists",
				err,
			)
		case 1452:
			// delete/update of a parent is blocked because children reference it
			return apperr.NewNotFound(
				"one or more referenced entity does not exist",
				err,
			)
		case 1451:
			// insert/update references a nonexistent parent
			return apperr.NewInvalid(
				"entity is still referenced by another entity",
				err,
			)
		}
	}

	return fmt.Errorf("%s: %w", operation, err)
}
