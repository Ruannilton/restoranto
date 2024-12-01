package infrastructure

import "github.com/Ruannilton/notification-service/domain/models"

type SMSSender struct {
}

func NewSmsSender() SMSSender {
	return SMSSender{}
}

func (client SMSSender) Send(e models.SMS) error {
	return nil
}
