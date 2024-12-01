package usecase

import (
	"errors"

	"github.com/Ruannilton/event-handler-service/domain/dependencies"
	"github.com/Ruannilton/go-msg-contracts/pkg/events"
	messagecontracts "github.com/Ruannilton/go-msg-contracts/pkg/message_contracts"
	"github.com/Ruannilton/go-msg-contracts/pkg/queues"
)

type CostumerCreatedUseCase struct {
	publisher dependencies.IMessagePublisher
}

func NewCostumerCreatedUseCase(publisher dependencies.IMessagePublisher) CostumerCreatedUseCase {
	return CostumerCreatedUseCase{
		publisher: publisher,
	}
}

func (useCase CostumerCreatedUseCase) Execute(event events.CostumerCreatedEvent) error {
	verifyEmail := messagecontracts.VerifyEmailMessage{
		Email:        event.Email,
		CostumerName: event.Name,
		NewCostumer:  true,
	}

	verifySms := messagecontracts.VerifySMSMessage{
		Phone:        event.Phone,
		CostumerName: event.Name,
		NewCostumer:  true,
	}

	emailValidationErr := useCase.publisher.Publish(verifyEmail, queues.EmailValidationQueue)
	smsValidationErr := useCase.publisher.Publish(verifySms, queues.SMSValidationQueue)

	err := errors.Join(emailValidationErr, smsValidationErr)

	if err != nil {
		return err
	}

	return nil
}
