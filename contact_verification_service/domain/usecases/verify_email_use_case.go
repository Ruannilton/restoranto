package usecases

import (
	"github.com/Ruannilton/contact-verification-service/domain/dependencies"
	"github.com/Ruannilton/contact-verification-service/domain/services"
	messagecontracts "github.com/Ruannilton/go-msg-contracts/pkg/message_contracts"
	"github.com/Ruannilton/go-msg-contracts/pkg/queues"
)

type VerifyEmailUseCase struct {
	publisher dependencies.IMessagePublisher
}

func NewVerifyEmailUseCase(publisher dependencies.IMessagePublisher) VerifyEmailUseCase {
	return VerifyEmailUseCase{
		publisher: publisher,
	}
}

func (useCase VerifyEmailUseCase) Execute(msg messagecontracts.VerifyEmailMessage) error {
	emailCode := services.GenerateEmailCode()

	emailTemplate := services.CreateVerifyEmailContent(msg, emailCode)

	sendMessage := messagecontracts.SendEmailMessage{
		Receiver: msg.Email,
		Subject:  useCase.getEmailSubject(msg),
		Body:     emailTemplate,
	}

	err := useCase.publisher.Publish(sendMessage, queues.SendEmailQueue)

	if err != nil {
		//TODO: handle err
		return err
	}

	return nil
}

func (useCase VerifyEmailUseCase) getEmailSubject(msg messagecontracts.VerifyEmailMessage) string {
	if msg.NewCostumer {
		return "Account creation"
	}
	return "Email address upate"
}
