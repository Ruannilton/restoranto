package usecases

import (
	"github.com/Ruannilton/notification-service/domain/dependencies"
	"github.com/Ruannilton/notification-service/domain/models"
)

type SendSMSUseCase struct {
	sender dependencies.ISMSSender
}

func NewSendSMSUseCase(sender dependencies.ISMSSender) SendSMSUseCase {
	return SendSMSUseCase{
		sender: sender,
	}
}

func (useCase SendSMSUseCase) Execute(sms models.SMS) error {
	err := useCase.sender.Send(sms)

	if err != nil {
		// TODO: handle error
		return err
	}

	return nil
}
