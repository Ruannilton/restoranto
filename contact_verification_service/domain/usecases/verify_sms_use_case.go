package usecases

import (
	"github.com/Ruannilton/contact-verification-service/domain/dependencies"
	"github.com/Ruannilton/contact-verification-service/domain/services"
	messagecontracts "github.com/Ruannilton/go-msg-contracts/pkg/message_contracts"
	"github.com/Ruannilton/go-msg-contracts/pkg/queues"
)

type VerifySMSUseCase struct {
	publisher dependencies.IMessagePublisher
}

func NewVerifySMSUseCase(publisher dependencies.IMessagePublisher) VerifySMSUseCase {
	return VerifySMSUseCase{
		publisher: publisher,
	}
}

func (useCase VerifySMSUseCase) Execute(msg messagecontracts.VerifySMSMessage) error {
	emailCode := services.GenerateEmailCode()

	emailTemplate := services.CreateVerifySMSContent(msg, emailCode)

	sendMessage := messagecontracts.SendSMSMessage{
		Number: msg.Phone,
		Body:   emailTemplate,
	}

	err := useCase.publisher.Publish(sendMessage, queues.SendEmailQueue)

	if err != nil {
		//TODO: handle err
		return err
	}

	return nil
}
