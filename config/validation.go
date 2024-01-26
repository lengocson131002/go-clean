package config

import (
	"sync"

	"github.com/lengocson131002/go-clean/pkg/validation"
)

var val validation.Validator
var onceVal sync.Once

func GetValidator() validation.Validator {
	if val == nil {
		onceVal.Do(func() {
			val = validation.NewGpValidator()
		})
	}
	return val
}
