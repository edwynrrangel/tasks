package tasks

import (
	"database/sql"
	"net/http"

	"github.com/edwynrrangel/tasks/errors"
)

func handlerError(err error, details errors.Details) *errors.Error {
	switch err {
	case sql.ErrNoRows:
		return errors.New(http.StatusNotFound, errors.MsgNotFound, details)
	case errNotFound:
		return errors.New(http.StatusNotFound, errMsgTaskNotFound, details)
	default:
		return errors.New(http.StatusInternalServerError, errors.MsgInternalServerError, details)
	}
}
