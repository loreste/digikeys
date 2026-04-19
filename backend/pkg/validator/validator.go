package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate checks struct tags and returns a human-readable error list.
func Validate(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	var messages []string
	for _, e := range validationErrors {
		messages = append(messages, formatFieldError(e))
	}

	return &ValidationError{Messages: messages}
}

type ValidationError struct {
	Messages []string
}

func (e *ValidationError) Error() string {
	return strings.Join(e.Messages, "; ")
}

func formatFieldError(e validator.FieldError) string {
	field := e.Field()

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("Le champ '%s' est obligatoire", field)
	case "email":
		return fmt.Sprintf("Le champ '%s' doit être un email valide", field)
	case "min":
		return fmt.Sprintf("Le champ '%s' doit avoir au moins %s caractères", field, e.Param())
	case "max":
		return fmt.Sprintf("Le champ '%s' ne doit pas dépasser %s caractères", field, e.Param())
	case "oneof":
		return fmt.Sprintf("Le champ '%s' doit être l'une des valeurs: %s", field, e.Param())
	default:
		return fmt.Sprintf("Le champ '%s' est invalide (%s)", field, e.Tag())
	}
}
