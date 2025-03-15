package interpreter

import "strconv"

// toFloat64 converts a value of any type to a float64 if possible.
// It returns the converted float64 value and a boolean indicating success.
func toFloat64(val any) (float64, bool) {
	if num, ok := val.(float64); ok {
		return num, ok
	}
	if str, ok := val.(string); ok {
		num, err := strconv.ParseFloat(str, 64)
		return num, err == nil
	}
	return 0, false
}

// isTruthy determines the truthiness of a value.
// It returns false for nil and false boolean values, and true for all other values.
func isTruthy(object any) bool {
	if object == nil {
		return false
	}
	switch v := object.(type) {
	case bool:
		return v
	default:
		return true
	}
}
