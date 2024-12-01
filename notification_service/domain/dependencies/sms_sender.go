package dependencies

import "github.com/Ruannilton/notification-service/domain/models"

type ISMSSender interface {
	Send(e models.SMS) error
}
