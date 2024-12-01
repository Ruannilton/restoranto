package presentation

import (
	"os"

	messagecontracts "github.com/Ruannilton/go-msg-contracts/pkg/message_contracts"
	"github.com/Ruannilton/go-msg-contracts/pkg/queues"
	"github.com/Ruannilton/notification-service/domain/models"
	usecases "github.com/Ruannilton/notification-service/domain/use_cases"
	amqp "github.com/rabbitmq/amqp091-go"
)

const EmailListenerName = "notification-service-email-listener"

type EmailListener struct {
	mqConnection *amqp.Connection
	emailSender  usecases.SendEmailUseCase
}

func NewEmailListener(mqConnection *amqp.Connection, emailSender usecases.SendEmailUseCase) EmailListener {
	return EmailListener{
		mqConnection: mqConnection,
		emailSender:  emailSender,
	}
}

func (listener EmailListener) ReceiveMessages(stopSig chan os.Signal) {
	channel, err := listener.mqConnection.Channel()

	if err != nil {
		//TODO: handle error
		return
	}

	defer channel.Close()

	messageChannel, err := channel.Consume(queues.SendEmailQueue, EmailListenerName, false, false, false, false, nil)

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
				message, err := messagecontracts.SendEmailMessage{}.FromBinary(msg.Body)

				if err != nil {
					msg.Nack(false, false)
					continue
				}

				email := models.Email{
					Receiver: message.Receiver,
					Subject:  message.Subject,
					Body:     message.Body,
				}

				err = listener.emailSender.Execute(email)

				if err != nil {
					msg.Nack(false, true)
					continue
				}

				msg.Ack(false)
			}
		}
	}()
}
