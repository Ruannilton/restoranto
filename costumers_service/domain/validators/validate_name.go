package validators

import (
	"costumers-api/domain/domain_errors"
	"strings"
)

func IsNameValid(name string) error {
	if strings.TrimSpace(name) == "" {
		return domain_errors.ErrInvalidName
	}
	return nil
}
