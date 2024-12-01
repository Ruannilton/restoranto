package validators

import (
	"costumers-api/domain/domain_errors"
	"regexp"
	"strconv"
	"strings"
)

func IsCpfValid(cpf string) error {

	if !isValidCPF(cpf) {
		return domain_errors.ErrInvalidCpf
	}

	return nil
}

func isValidCPF(cpf string) bool {

	re := regexp.MustCompile(`[^\d]`)
	cpf = re.ReplaceAllString(cpf, "")

	if len(cpf) != 11 || strings.Repeat(string(cpf[0]), 11) == cpf {
		return false
	}

	if !validateDigit(cpf, 9) {
		return false
	}

	return validateDigit(cpf, 10)
}

func validateDigit(cpf string, length int) bool {
	sum := 0
	for i := 0; i < length; i++ {
		digit, _ := strconv.Atoi(string(cpf[i]))
		sum += digit * (length + 1 - i)
	}
	remainder := sum % 11
	checkDigit := 0
	if remainder >= 2 {
		checkDigit = 11 - remainder
	}
	return strconv.Itoa(checkDigit) == string(cpf[length])
}
