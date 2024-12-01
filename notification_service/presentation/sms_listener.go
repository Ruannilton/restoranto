package presentation

import (
	"os"

	messagecontracts "github.com/Ruannilton/go-msg-contracts/pkg/message_contracts"
	"github.com/Ruannilton/go-msg-contracts/pkg/queues"
	"github.com/Ruannilton/notification-service/domain/models"
	usecases "github.com/Ruannilton/notification-service/domain/use_cases"
	amqp "github.com/rabbitmq/amqp091-go"
)

type SMSListener struct {
	mqConnection *amqp.Connection
	smsSender    usecases.SendSMSUseCase
}

const SMSListenerName = "notification-service-sms-listener"

func NewSMSListener(mqConnection *amqp.Connection, smsSender usecases.SendSMSUseCase) SMSListener {
	return SMSListener{
		mqConnection: mqConnection,
		smsSender:    smsSender,
	}
}

func (listener SMSListener) ReceiveMessages(stopSig chan os.Signal) {
	channel, err := listener.mqConnection.Channel()

	if err != nil {
		//TODO: handle error
		return
	}

	defer channel.Close()

	messageChannel, err := channel.Consume(queues.SMSValidationQueue, SMSListenerName, false, false, false, false, nil)

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
				message, err := messagecontracts.SendSMSMessage{}.FromBinary(msg.Body)

				if err != nil {
					msg.Nack(false, false)
					continue
				}

				sms := models.SMS{
					Receiver: message.Number,
					Body:     message.Body,
				}

				err = listener.smsSender.Execute(sms)

				if err != nil {
					msg.Nack(false, true)
					continue
				}

				msg.Ack(false)
			}
		}
	}()
}
