package auth

import (
	"net/http"

	"github.com/edwynrrangel/tasks/errors"
)

func handlerError(err error, details errors.Details) *errors.Error {
	switch err {
	case errFirstLogin:
		return errors.New(http.StatusAccepted, errMsgFirstLogin, details)
	case errUnauthorized:
		return errors.New(http.StatusUnauthorized, errMsgUnauthorized, details)
	case errInactiveUser:
		return errors.New(http.StatusUnauthorized, errMsgInactiveUser, details)
	default:
		return errors.New(http.StatusInternalServerError, errors.MsgInternalServerError, details)
	}
}
