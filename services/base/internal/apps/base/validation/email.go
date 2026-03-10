package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/sdq-codes/usegro-api/internal/apps/base/dto"
	"github.com/sdq-codes/usegro-api/internal/interface/validation"
	"github.com/sdq-codes/usegro-api/pkg/exception"
)

type EmailValidation struct{}

func (Ev *EmailValidation) EmailVerificationValidation(dti dto.EmailVerificationDTI) error {
	v, _ := validation.GetValidator()
	if err := v.Struct(dti); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			return exception.NewValidationFailedErrors(validationErrs)
		}
	}
	return nil
}
