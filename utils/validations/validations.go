package validations

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func AddCustomValidations() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("depends_on", dependsOn)

		if err != nil {
			return err
		}
	}
	return nil
}
