package auth

import "errors"

const (
	errMsgInvalidCredentials string = "credenciales no válidas"
	errMsgInvalidToken       string = "token no válido"
	errMsgFirstLogin         string = "debe cambiar su contraseña inicial"
	errMsgUnauthorized       string = "credenciales no válidas"
	errMsgInactiveUser       string = "usuario inactivo"
)

var (
	errFirstLogin   error = errors.New(errMsgFirstLogin)
	errUnauthorized error = errors.New(errMsgUnauthorized)
	errInactiveUser error = errors.New(errMsgInactiveUser)
)
