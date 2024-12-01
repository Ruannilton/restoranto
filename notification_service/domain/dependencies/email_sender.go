package dependencies

import "github.com/Ruannilton/notification-service/domain/models"

type IEmailSender interface {
	Send(e models.Email) error
}
