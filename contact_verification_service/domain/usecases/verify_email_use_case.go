package usecases

import (
	"fmt"
	"time"

	"github.com/Ruannilton/contact-verification-service/domain/dependencies"
	"github.com/Ruannilton/contact-verification-service/domain/services"
	messagecontracts "github.com/Ruannilton/go-msg-contracts/pkg/message_contracts"
	"github.com/Ruannilton/go-msg-contracts/pkg/queues"
)

type VerifyEmailUseCase struct {
	publisher dependencies.IMessagePublisher
	cache     dependencies.IDistributedCacheRepository
}

func NewVerifyEmailUseCase(publisher dependencies.IMessagePublisher, cache dependencies.IDistributedCacheRepository) VerifyEmailUseCase {
	return VerifyEmailUseCase{
		publisher: publisher,
		cache:     cache,
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

	key := fmt.Sprintf("validate_email:%s", msg.Email)
	cache_err := useCase.cache.Set(key, emailCode, time.Hour)

	if cache_err != nil {
		//TODO: handle err
		return cache_err
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
