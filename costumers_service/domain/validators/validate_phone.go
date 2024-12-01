package validators

import (
	"costumers-api/domain/domain_errors"

	"github.com/nyaruka/phonenumbers"
)

func IsPhoneValid(phone string) error {
	parsedNumber, err := phonenumbers.Parse(phone, "")
	if err != nil {
		return domain_errors.ErrInvalidPhone
	}

	if !phonenumbers.IsValidNumber(parsedNumber) {
		return domain_errors.ErrInvalidPhone
	}

	return nil
}
