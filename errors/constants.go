package errors

import "errors"

const (
	MsgInternalServerError string = "Error interno del servidor"
	MsgInvalidBody         string = "Cuerpo de la petición no válido"
	MsgInvalidQuery        string = "Query de la petición no válido"
	MsgNotFound            string = "Registro no encontrado"
)

var (
	ErrInternalServerError error = errors.New(MsgInternalServerError)
	ErrInvalidBody         error = errors.New(MsgInvalidBody)
	ErrInvalidQuery        error = errors.New(MsgInvalidQuery)
)
