package validation

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/usegro/services/crm/internal/apps/crm/dto"
	"github.com/usegro/services/crm/internal/interface/validation"
	"github.com/usegro/services/crm/pkg/exception"
)

// CreateTag validation
func CreateTagValidation(dti dto.TagCreateDTO) error {
	v, _ := validation.GetValidator()
	if err := v.Struct(dti); err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			return exception.NewValidationFailedErrors(validationErrs)
		}
	}
	return nil
}

// UpdateTagName validation
func UpdateTagNameValidation(dti dto.TagUpdateDTO) error {
	v, _ := validation.GetValidator()
	if err := v.Struct(dti); err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			return exception.NewValidationFailedErrors(validationErrs)
		}
	}
	return nil
}
