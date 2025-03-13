package scanner

import (
	"fmt"
	"strconv"
)

func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphaNumeric(c rune) bool {
	return isAlpha(c) || isDigit(c)
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

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
