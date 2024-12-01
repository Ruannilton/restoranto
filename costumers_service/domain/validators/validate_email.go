package validators

import (
	"costumers-api/domain/domain_errors"
	"net/mail"
)

func IsEmailValid(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return domain_errors.ErrInvalidEmail
	}

	return nil
}
