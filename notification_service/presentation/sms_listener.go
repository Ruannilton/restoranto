package presentation

import (
	"os"

	messagecontracts "github.com/Ruannilton/go-msg-contracts/pkg/message_contracts"
	"github.com/Ruannilton/go-msg-contracts/pkg/queues"
	"github.com/Ruannilton/notification-service/domain/dependencies"
	"github.com/Ruannilton/notification-service/domain/models"
	usecases "github.com/Ruannilton/notification-service/domain/use_cases"
)

type SMSListener struct {
	listener  dependencies.IMessageListener
	smsSender usecases.SendSMSUseCase
}

const SMSListenerName = "notification-service-sms-listener"

func NewSMSListener(listener dependencies.IMessageListener, smsSender usecases.SendSMSUseCase) SMSListener {
	return SMSListener{
		listener:  listener,
		smsSender: smsSender,
	}
}

func (listener SMSListener) ReceiveMessages(stopSig chan os.Signal) {
	defer listener.listener.Close()

	messageChannel, err := listener.listener.Listen(queues.EmailValidationQueue)

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
				message, err := messagecontracts.SendSMSMessage{}.FromBinary(msg.GetContent())

				if err != nil {
					msg.Nack(false)
					continue
				}

				sms := models.SMS{
					Receiver: message.Number,
					Body:     message.Body,
				}

				err = listener.smsSender.Execute(sms)

				if err != nil {
					msg.Nack(true)
					continue
				}

				msg.Ack()
			}
		}
	}()
}
