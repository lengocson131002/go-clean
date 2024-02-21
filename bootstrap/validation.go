package bootstrap

import (
	"github.com/lengocson131002/go-clean/pkg/validation"
)

func GetValidator() validation.Validator {
	return validation.NewGpValidator()
}
