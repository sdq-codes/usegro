package validation

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/usegro/services/crm/internal/apps/crm/dto"
	"github.com/usegro/services/crm/internal/interface/validation"
	"github.com/usegro/services/crm/pkg/exception"
)

func CRMUserOrganizationCreateValidation(dti dto.CrmUserOrganizationDTI) error {
	v, _ := validation.GetValidator()
	if err := v.Struct(dti); err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			return exception.NewValidationFailedErrors(validationErrs)
		}
	}
	return nil
}

func CRMUserOrganizationUpdateValidation(dti dto.CrmUserOrganizationDTI) error {
	v, _ := validation.GetValidator()
	if err := v.Struct(dti); err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			return exception.NewValidationFailedErrors(validationErrs)
		}
	}
	return nil
}

func CRMUserOrganizationSalesChannelTypeValidation(dti dto.CrmUserOrganizationSalesChannelTypeDTI) error {
	v, _ := validation.GetValidator()
	if err := v.Struct(dti); err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			return exception.NewValidationFailedErrors(validationErrs)
		}
	}
	return nil
}

func CRMUserOrganizationStockProductTypeValidation(dti dto.CrmUserOrganizationStockProductTypeDTI) error {
	v, _ := validation.GetValidator()
	if err := v.Struct(dti); err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			return exception.NewValidationFailedErrors(validationErrs)
		}
	}
	return nil
}
