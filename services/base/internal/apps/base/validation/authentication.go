package validation

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/sdq-codes/usegro-api/internal/apps/base/dto"
	"github.com/sdq-codes/usegro-api/internal/interface/validation"
	"github.com/sdq-codes/usegro-api/pkg/exception"
	"regexp"
)

type AuthenticationValidation struct{}

func passwordHasLowercase(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`[a-z]`).MatchString(fl.Field().String())
}

func passwordHasUppercase(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`[A-Z]`).MatchString(fl.Field().String())
}

func passwordHasDigit(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`\d`).MatchString(fl.Field().String())
}

func passwordHasSpecialChar(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};':"\\|,.<>\/?]`).MatchString(fl.Field().String())
}

func (av *AuthenticationValidation) CredentialsValidation(dti dto.RegisterUserDTI) error {
	v, _ := validation.GetValidator()
	if err := v.RegisterValidation("passwd_lower", passwordHasLowercase); err != nil {
		return err
	}
	if err := v.RegisterValidation("passwd_upper", passwordHasUppercase); err != nil {
		return err
	}
	if err := v.RegisterValidation("passwd_digit", passwordHasDigit); err != nil {
		return err
	}
	if err := v.RegisterValidation("passwd_special", passwordHasSpecialChar); err != nil {
		return err
	}

	if err := v.Struct(dti); err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			// Map custom error messages
			customMessages := map[string]string{
				"RegisterUserDTI.Password.passwd_lower":   "Password must contain at least one lowercase letter",
				"RegisterUserDTI.Password.passwd_upper":   "Password must contain at least one uppercase letter",
				"RegisterUserDTI.Password.passwd_digit":   "Password must contain at least one digit",
				"RegisterUserDTI.Password.passwd_special": "Password must contain at least one special character",
			}

			for i, fieldErr := range validationErrs {
				key := fieldErr.StructNamespace() + "." + fieldErr.Tag()
				if msg, ok := customMessages[key]; ok {
					validationErrs[i] = &validation.CustomFieldError{FieldError: fieldErr, Message: msg}
				}
			}
			return exception.NewValidationFailedErrors(validationErrs)
		}
		return err
	}

	return nil
}

func (av *AuthenticationValidation) ForgotPasswordEmailValidation(dti dto.ForgotPasswordDTI) error {
	v, _ := validation.GetValidator()
	if err := v.Struct(dti); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			return exception.NewValidationFailedErrors(validationErrs)
		}
	}
	return nil
}

func (av *AuthenticationValidation) UserExistValidation(dti dto.UserExistDTI) error {
	v, _ := validation.GetValidator()
	if err := v.Struct(dti); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			return exception.NewValidationFailedErrors(validationErrs)
		}
	}
	return nil
}

func (av *AuthenticationValidation) ResetPasswordEmailValidation(dti dto.ResetPasswordDTI) error {
	v, _ := validation.GetValidator()
	if err := v.RegisterValidation("passwd_lower", passwordHasLowercase); err != nil {
		return err
	}
	if err := v.RegisterValidation("passwd_upper", passwordHasUppercase); err != nil {
		return err
	}
	if err := v.RegisterValidation("passwd_digit", passwordHasDigit); err != nil {
		return err
	}
	if err := v.RegisterValidation("passwd_special", passwordHasSpecialChar); err != nil {
		return err
	}

	if err := v.Struct(dti); err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			// Map custom error messages
			customMessages := map[string]string{
				"RegisterUserDTI.Password.passwd_lower":   "Password must contain at least one lowercase letter",
				"RegisterUserDTI.Password.passwd_upper":   "Password must contain at least one uppercase letter",
				"RegisterUserDTI.Password.passwd_digit":   "Password must contain at least one digit",
				"RegisterUserDTI.Password.passwd_special": "Password must contain at least one special character",
			}

			for i, fieldErr := range validationErrs {
				key := fieldErr.StructNamespace() + "." + fieldErr.Tag()
				if msg, ok := customMessages[key]; ok {
					validationErrs[i] = &validation.CustomFieldError{FieldError: fieldErr, Message: msg}
				}
			}
			return exception.NewValidationFailedErrors(validationErrs)
		}
		return err
	}

	return nil
}

func (av *AuthenticationValidation) RequestEmailCodeValidation(dti dto.RequestEmailCodeDTI) error {
	v, _ := validation.GetValidator()
	if err := v.Struct(dti); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			return exception.NewValidationFailedErrors(validationErrs)
		}
	}
	return nil
}

func (av *AuthenticationValidation) VerifyEmailCodeValidation(dti dto.VerifyEmailCodeDTI) error {
	v, _ := validation.GetValidator()
	if err := v.Struct(dti); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			return exception.NewValidationFailedErrors(validationErrs)
		}
	}
	return nil
}
