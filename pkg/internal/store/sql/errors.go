package sql

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/michaljemala/pqerror"

	"github.com/michaljemala/payments-sample/pkg/internal/errors"
)

func WrapSelectError(err error, msg string) error {
	switch err {
	case nil:
		return nil
	case sql.ErrNoRows:
		return errors.Generic(errors.ErrCodeGenericNotFound, msg, err.Error())
	}
	return errors.DataAccess(errors.ErrCodeDataAccessSelectFailed, msg, err.Error())
}

func WrapInsertError(err error, msg string) error {
	switch err := err.(type) {
	case nil:
		return nil
	case *pq.Error: //must unwrap the errors to extract specific DB codes
		switch err.Code {
		case pqerror.UniqueViolation:
			return errors.Generic(errors.ErrCodeGenericAlreadyExists, msg, err.Error())
		}
	}
	return errors.DataAccess(errors.ErrCodeDataAccessInsertFailed, msg, err.Error())
}

func WrapDeleteError(err error, msg string) error {
	switch err {
	case nil:
		return nil
	case sql.ErrNoRows:
		return errors.Generic(errors.ErrCodeGenericNotFound, msg, err.Error())
	}
	return errors.DataAccess(errors.ErrCodeDataAccessDeleteFailed, msg, err.Error())
}

func WrapUpdateError(err error, msg string) error {
	switch err {
	case nil:
		return nil
	case sql.ErrNoRows:
		return errors.Generic(errors.ErrCodeGenericNotFound, msg, err.Error())
	}
	return errors.DataAccess(errors.ErrCodeDataAccessUpdateFailed, msg, err.Error())
}
