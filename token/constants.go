package token

import "errors"

const (
	errMsgNotFound              string = "token no encontrado"
	errMsgFirstLogin            string = "debe cambiar su contraseña inicial"
	errMsgExpiredInvalidMissing string = "Tu sesión ha expirado, es inválida o no existe"
)

var (
	errNotFound error = errors.New(errMsgNotFound)
)
