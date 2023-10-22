package users

import (
	"net/http"

	"github.com/edwynrrangel/tasks/errors"
)

func handlerError(err error, details errors.Details) *errors.Error {
	switch err {
	case errUserAlreadyExists:
		return errors.New(http.StatusConflict, err.Error(), details)
	default:
		return errors.New(http.StatusInternalServerError, errors.MsgInternalServerError, details)
	}
}
