package validator

import (
	"fmt"

	"github.com/edwynrrangel/tasks/errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var Validator = validator.New()

func BodyParser(ctx *fiber.Ctx, data interface{}) error {
	if err := ctx.BodyParser(data); err != nil {
		return err
	}
	if err := validateSchema(data); err != nil {
		return err
	}
	return nil
}

func QueryParser(ctx *fiber.Ctx, data interface{}) error {
	if err := ctx.QueryParser(data); err != nil {
		return err
	}
	if err := validateSchema(data); err != nil {
		return err
	}
	return nil
}

func validateSchema(data interface{}) error {
	if err := Validator.Struct(data); err != nil {
		var (
			ve  validator.ValidationErrors
			out string
		)
		if errors.As(err, &ve) {
			for i, fe := range ve {
				out += msgForTag(fe)
				if i != len(ve)-1 {
					out += ", "
				}
			}
		}

		return fmt.Errorf(out)
	}
	return nil
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("El campo %s es requerido.", fe.Field())
	case "eqfield":
		return fmt.Sprintf("El campo %s debe ser igual al campo %s.", fe.Field(), fe.Param())
	case "min":
		return fmt.Sprintf("El valor del campo %s debe ser mayor o igual a %s.", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("El valor del campo %s debe ser menor o igual a %s.", fe.Field(), fe.Param())
	case "oneof":
		return fmt.Sprintf("El valor del campo %s debe ser uno de los siguientes: %s.", fe.Field(), fe.Param())
	case "datetime":
		return fmt.Sprintf("El valor del campo %s debe ser una fecha con el formato %s.", fe.Field(), fe.Param())
	case "uuid4":
		return fmt.Sprintf("El valor del campo %s debe ser un UUID v4.", fe.Field())
	case "required_with":
		return fmt.Sprintf("El campo %s es requerido cuando el campo %s tiene un valor.", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("Fallo la etiqueta %s", fe.Tag())
	}
}
