package validators

import (
	"costumers-api/domain/domain_errors"
	"time"
)

const MinAgeInHours = 157.788

func IsBirthdateValid(birthdate time.Time) error {

	if time.Now().Before(birthdate) {
		return domain_errors.ErrInvalidBirthdate
	}

	if time.Now().AddDate(-18, 0, 0).Before(birthdate) {
		return domain_errors.ErrInvalidBirthdate
	}

	if time.Now().AddDate(-65, 0, 0).After(birthdate) {
		return domain_errors.ErrInvalidBirthdate
	}

	timeElapsed := time.Since(birthdate)

	hours := timeElapsed.Hours()

	if hours < MinAgeInHours {
		return domain_errors.ErrUnderage
	}

	return nil
}
