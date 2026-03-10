package validation

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/usegro/services/crm/pkg/exception"
)

type Validation struct{}

func (v *Validation) Validate(validate validator.Validate, dti interface{}) *exception.ExceptionErrors {
	if err := validate.Struct(dti); err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			return exception.NewValidationFailedErrors(validationErrs)
		}
	}
	return nil
}
