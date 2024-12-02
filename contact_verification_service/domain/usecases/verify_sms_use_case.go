package usecases

import (
	"fmt"
	"time"

	"github.com/Ruannilton/contact-verification-service/domain/dependencies"
	"github.com/Ruannilton/contact-verification-service/domain/services"
	messagecontracts "github.com/Ruannilton/go-msg-contracts/pkg/message_contracts"
	"github.com/Ruannilton/go-msg-contracts/pkg/queues"
)

type VerifySMSUseCase struct {
	publisher dependencies.IMessagePublisher
	cache     dependencies.IDistributedCacheRepository
}

func NewVerifySMSUseCase(publisher dependencies.IMessagePublisher, cache dependencies.IDistributedCacheRepository) VerifySMSUseCase {
	return VerifySMSUseCase{
		publisher: publisher,
		cache:     cache,
	}
}

func (useCase VerifySMSUseCase) Execute(msg messagecontracts.VerifySMSMessage) error {
	emailCode := services.GenerateEmailCode()

	emailTemplate := services.CreateVerifySMSContent(msg, emailCode)

	sendMessage := messagecontracts.SendSMSMessage{
		Number: msg.Phone,
		Body:   emailTemplate,
	}

	key := fmt.Sprintf("validate_phone:%s", msg.Phone)
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
