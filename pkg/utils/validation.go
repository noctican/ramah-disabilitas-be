package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string
		for _, e := range validationErrors {
			field := e.Field()
			switch field {
			case "ConfirmPassword":
				field = "Konfirmasi Password"
			case "Password":
				field = "Password"
			case "Email":
				field = "Email"
			case "Name":
				field = "Nama"
			}

			switch e.Tag() {
			case "required":
				errorMessages = append(errorMessages, fmt.Sprintf("%s wajib diisi", field))
			case "email":
				errorMessages = append(errorMessages, fmt.Sprintf("%s harus berupa email yang valid", field))
			case "min":
				errorMessages = append(errorMessages, fmt.Sprintf("%s minimal %s karakter", field, e.Param()))
			case "eqfield":
				errorMessages = append(errorMessages, fmt.Sprintf("%s harus sama dengan %s", field, e.Param()))
			default:
				errorMessages = append(errorMessages, fmt.Sprintf("%s tidak valid", field))
			}
		}
		return strings.Join(errorMessages, ", ")
	}
	return err.Error()
}
