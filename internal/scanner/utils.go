package scanner

import (
	"fmt"
	"strconv"
)

// isAlpha checks if a rune is an alphabetic character (a-z, A-Z, _).
func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

// isAlphaNumeric checks if a rune is an alphanumeric character (a-z, A-Z, 0-9, _).
func isAlphaNumeric(c rune) bool {
	return isAlpha(c) || isDigit(c)
}

// isDigit checks if a rune is a digit (0-9).
func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

// formatFloat formats a float number to a string.
// If the number is effectively an integer, it will be formatted with one decimal place.
func formatFloat(text string) (string, error) {
	num, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return "", err
	}

	// check if number is effectively an integer
	if num == float64(int(num)) {
		return fmt.Sprintf("%.1f", num), nil
	}
	return fmt.Sprintf("%g", num), nil
}
