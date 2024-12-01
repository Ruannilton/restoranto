package validators

import (
	"costumers-api/domain/domain_errors"
	"strconv"
	"strings"
)

func IsCnpjValid(cnpj string) error {
	// Remove any non-digit characters
	cnpj = strings.ReplaceAll(cnpj, ".", "")
	cnpj = strings.ReplaceAll(cnpj, "-", "")
	cnpj = strings.ReplaceAll(cnpj, "/", "")

	// CNPJ must be 14 digits
	if len(cnpj) != 14 {
		return domain_errors.ErrInvalidCnpj
	}

	// Check if all digits are the same, which is an invalid CNPJ
	if allSameDigits(cnpj) {
		return domain_errors.ErrInvalidCnpj
	}

	// Convert CNPJ string to a slice of integers
	cnpjDigits := make([]int, 14)
	for i := 0; i < 14; i++ {
		digit, err := strconv.Atoi(string(cnpj[i]))
		if err != nil {
			return domain_errors.ErrInvalidCnpj
		}
		cnpjDigits[i] = digit
	}

	// Validate first check digit
	if cnpjDigits[12] != calculateCheckDigit(cnpjDigits, 12) {
		return domain_errors.ErrInvalidCnpj
	}

	// Validate second check digit
	if cnpjDigits[13] != calculateCheckDigit(cnpjDigits, 13) {
		return domain_errors.ErrInvalidCnpj
	}

	return nil
}

func allSameDigits(cnpj string) bool {
	for i := 1; i < len(cnpj); i++ {
		if cnpj[i] != cnpj[0] {
			return false
		}
	}
	return true
}

// calculateCheckDigit calculates the CNPJ check digit for the given position
func calculateCheckDigit(cnpj []int, length int) int {
	var weights []int
	if length == 12 {
		weights = []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	} else if length == 13 {
		weights = []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	} else {
		return -1
	}

	sum := 0
	for i := 0; i < length; i++ {
		sum += cnpj[i] * weights[i]
	}

	remainder := sum % 11
	if remainder < 2 {
		return 0
	}
	return 11 - remainder
}
