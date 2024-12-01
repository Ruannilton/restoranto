package usecases

import (
	"github.com/Ruannilton/notification-service/domain/dependencies"
	"github.com/Ruannilton/notification-service/domain/models"
)

type SendEmailUseCase struct {
	sender dependencies.IEmailSender
}

func NewSendEmailUseCase(sender dependencies.IEmailSender) SendEmailUseCase {
	return SendEmailUseCase{
		sender: sender,
	}
}

func (useCase SendEmailUseCase) Execute(email models.Email) error {
	err := useCase.sender.Send(email)

	if err != nil {
		//TODO: handle error
		return err
	}

	return nil
}
