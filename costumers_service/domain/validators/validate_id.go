package validators

import "costumers-api/domain/domain_errors"

func IsIdValid(id int) error {
	if id <= 0 {
		return domain_errors.ErrInvalidId
	}
	return nil
}
