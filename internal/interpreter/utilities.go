package interpreter

import (
	"codecrafters-interpreter-go/internal/ast"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

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

// isEqual compares two values of any type and determines if they are deeply equal.
// It returns true if both values are nil or if they are deeply equal according to reflect.DeepEqual.
// If one value is nil and the other is not, it returns false.
//
// Parameters:
//   - a: The first value to compare.
//   - b: The second value to compare.
//
// Returns:
//   - bool: True if the values are deeply equal, false otherwise.
func isEqual(a any, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	return reflect.DeepEqual(a, b)
}

// throwRuntimeError logs a runtime error message to the standard error output
// along with the line number where the error occurred, and then terminates
// the program with an exit code of 70.
//
// Parameters:
//   - token: An ast.Token representing the source code token where the error occurred.
//     The token's Line field is used to indicate the line number.
//   - message: A string containing the error message to be displayed.
//
// This function is typically used to handle unrecoverable runtime errors in
// the interpreter.
func throwRuntimeError(token ast.Token, message string) {
	fmt.Fprintf(os.Stderr, "%s \n[line %d]\n", message, token.Line)
	os.Exit(70)
}
