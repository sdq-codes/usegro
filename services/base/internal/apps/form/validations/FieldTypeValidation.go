package validation

import "github.com/sdq-codes/usegro-api/internal/apps/form/validations/Rules"

type ValidationRule struct {
	Key     string
	Rule    string
	Param   string
	Message string
}

type FieldValidator struct {
	FieldName   string
	FieldType   string
	Validations []ValidationRule
}

// Validate runs all validations for a field
func (fv FieldValidator) Validate(value string) Rules.ValidationResult {
	result := Rules.ValidationResult{Valid: true, Errors: []string{}}

	for _, v := range fv.Validations {
		ok, errMsg := Rules.ApplyRule(v.Rule, value, v.Param, v.Message)
		if !ok {
			result.Valid = false
			result.Errors = append(result.Errors, errMsg)
		}
	}

	return result
}

// ValidateForm validates multiple fields at once
func ValidateForm(fields []FieldValidator, values map[string]string) map[string]Rules.ValidationResult {
	results := make(map[string]Rules.ValidationResult)

	for _, field := range fields {
		val := values[field.FieldName]
		results[field.FieldName] = field.Validate(val)
	}

	return results
}
