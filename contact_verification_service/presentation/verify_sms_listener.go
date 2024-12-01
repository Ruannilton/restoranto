package presentation

import (
	"os"

	"github.com/Ruannilton/contact-verification-service/domain/usecases"
	messagecontracts "github.com/Ruannilton/go-msg-contracts/pkg/message_contracts"
	"github.com/Ruannilton/go-msg-contracts/pkg/queues"
	amqp "github.com/rabbitmq/amqp091-go"
)

const VerifySMSListenerName = "contact-verification-verify-sms-listener"

type VerifySMSListener struct {
	mqConnection *amqp.Connection
	useCase      usecases.VerifySMSUseCase
}

func NewVerifySMSListener(mqConnection *amqp.Connection, useCase usecases.VerifySMSUseCase) VerifySMSListener {
	return VerifySMSListener{
		mqConnection: mqConnection,
		useCase:      useCase,
	}
}

func (listener VerifySMSListener) ReceiveMessages(stopSig chan os.Signal) {
	channel, err := listener.mqConnection.Channel()

	if err != nil {
		//TODO: handle error
		return
	}

	defer channel.Close()

	messageChannel, err := channel.Consume(queues.SMSValidationQueue, VerifySMSListenerName, false, false, false, false, nil)

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
				message, err := messagecontracts.VerifySMSMessage{}.FromBinary(msg.Body)

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
