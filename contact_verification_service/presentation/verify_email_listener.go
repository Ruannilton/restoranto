package presentation

import (
	"os"

	"github.com/Ruannilton/contact-verification-service/domain/usecases"
	messagecontracts "github.com/Ruannilton/go-msg-contracts/pkg/message_contracts"
	"github.com/Ruannilton/go-msg-contracts/pkg/queues"
	amqp "github.com/rabbitmq/amqp091-go"
)

const VerifyEmailListenerName = "contact-verification-verify-email-listener"

type VerifyEmailListener struct {
	mqConnection *amqp.Connection
	useCase      usecases.VerifyEmailUseCase
}

func NewVerifyEmailListener(mqConnection *amqp.Connection, useCase usecases.VerifyEmailUseCase) VerifyEmailListener {
	return VerifyEmailListener{
		mqConnection: mqConnection,
		useCase:      useCase,
	}
}

func (listener VerifyEmailListener) ReceiveMessages(stopSig chan os.Signal) {
	channel, err := listener.mqConnection.Channel()

	if err != nil {
		//TODO: handle error
		return
	}

	defer channel.Close()

	messageChannel, err := channel.Consume(queues.EmailValidationQueue, VerifyEmailListenerName, false, false, false, false, nil)

	if err != nil {
		//TODO: handle error
		return
	}

	go func() {
		for {
			select {
			case <-stopSig:
				return
			case msg := <-messageChannel:
				message, err := messagecontracts.VerifyEmailMessage{}.FromBinary(msg.Body)

				if err != nil {
					//TODO: handle error
					msg.Nack(false, false)
					continue
				}

				err = listener.useCase.Execute(message)

				if err != nil {
					//TODO: handle error
					msg.Nack(false, true)
					continue
				}

				msg.Ack(false)
			}
		}
	}()
}
