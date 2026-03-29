package helpers

import (
	"fmt"
	"github.com/sdq-codes/usegro-api/internal/apps/form/models"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

// ValidateSubmissionAnswers validates answers against form fields using "validations"
func ValidateSubmissionAnswers(fields []models.FormVersionField, answers map[string]interface{}) error {
	fieldMap := make(map[string]models.FormVersionField)
	for _, f := range fields {
		fieldId := f.ID
		fieldMap[fieldId] = f
	}

	// 1. Required fields
	for _, f := range fields {
		fieldId := f.ID

		if f.Required {
			if _, exists := answers[fieldId]; !exists {
				return fmt.Errorf("%s is required", f.Label)
			}
		}
	}

	// 2. No extra fields
	for key := range answers {
		if _, ok := fieldMap[key]; !ok {
			return fmt.Errorf("unexpected field in answers: %s", key)
		}
	}

	for key, val := range answers {
		f := fieldMap[key]

		for _, rule := range f.Validations {
			key, _ := rule["key"]
			message, _ := rule["message"]

			switch key {
			case "required":
				if val == nil || val == "" {
					return fmt.Errorf(message)
				}

			case "minLength":
				if strVal, ok := val.(string); ok {
					ruleInt, err := strconv.Atoi(rule["rule"])
					if err != nil {
						return fmt.Errorf(message)
					}
					if len(strVal) < ruleInt {
						return fmt.Errorf(message)
					}
				}

			case "maxLength":
				if strVal, ok := val.(string); ok {
					ruleInt, err := strconv.Atoi(rule["rule"])
					if err != nil {
						return fmt.Errorf(message)
					}
					if len(strVal) > ruleInt {
						return fmt.Errorf(message)
					}
				}

			case "min":
				ruleInt, err := strconv.Atoi(rule["rule"])
				if err != nil {
					return fmt.Errorf(message)
				}
				if num := val.(int); num < ruleInt {
					return fmt.Errorf(message)
				}

			case "max":
				ruleInt, err := strconv.Atoi(rule["rule"])
				if err != nil {
					return fmt.Errorf(message)
				}
				if num := val.(int); num > ruleInt {
					return fmt.Errorf(message)
				}

			case "email":
				if strVal, ok := val.(string); ok {
					re := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
					if !re.MatchString(strVal) {
						return fmt.Errorf(message)
					}
				} else {
					return fmt.Errorf("%s must be a string (email)", f.Label)
				}

			case "date":
				if strVal, ok := val.(string); ok {
					if _, err := time.Parse("2006-01-02", strVal); err != nil {
						return fmt.Errorf(message)
					}
				} else {
					return fmt.Errorf("%s must be a string (date)", f.Label)
				}

			case "boolean":
				if _, ok := val.(bool); !ok {
					return fmt.Errorf(message)
				}

			case "array":
				if reflect.TypeOf(val).Kind() != reflect.Slice {
					return fmt.Errorf(message)
				}
			}
		}
	}

	return nil
}

// helper: convert any numeric type to float64
func toFloat(v interface{}) float64 {
	switch n := v.(type) {
	case int:
		return float64(n)
	case int32:
		return float64(n)
	case int64:
		return float64(n)
	case float32:
		return float64(n)
	case float64:
		return n
	default:
		return 0
	}
}

func mergeMaps(a, b map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range a {
		result[k] = v
	}
	for k, v := range b {
		result[k] = v
	}
	return result
}
