package Rules

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// ValidationResult holds the result of a validation
type ValidationResult struct {
	Valid  bool
	Errors []string
}

// applyRule applies a single validation rule
func ApplyRule(rule string, value string, param string, message string) (bool, string) {
	switch rule {
	case "required":
		if value == "" {
			return false, message
		}
	case "minLength":
		min, _ := strconv.Atoi(param)
		if len(value) < min {
			return false, fmt.Sprintf("%s (minimum %d characters)", message, min)
		}
	case "maxLength":
		max, _ := strconv.Atoi(param)
		if len(value) > max {
			return false, fmt.Sprintf("%s (maximum %d characters)", message, max)
		}
	case "regex":
		re, err := regexp.Compile(param)
		if err == nil && !re.MatchString(value) {
			return false, message
		}
	case "number":
		_, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return false, message
		}
	case "date":
		_, err := time.Parse("2006-01-02", value)
		if err != nil {
			return false, message
		}
	case "time":
		_, err := time.Parse("15:04", value)
		if err != nil {
			return false, message
		}
	case "datetime":
		_, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return false, message
		}
	default:
		return true, ""
	}
	return true, ""
}
