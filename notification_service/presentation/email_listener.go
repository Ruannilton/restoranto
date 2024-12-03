package presentation

import (
	"os"

	messagecontracts "github.com/Ruannilton/go-msg-contracts/pkg/message_contracts"
	"github.com/Ruannilton/go-msg-contracts/pkg/queues"
	"github.com/Ruannilton/notification-service/domain/dependencies"
	"github.com/Ruannilton/notification-service/domain/models"
	usecases "github.com/Ruannilton/notification-service/domain/use_cases"
)

const EmailListenerName = "notification-service-email-listener"

type EmailListener struct {
	emailSender usecases.SendEmailUseCase
	listener    dependencies.IMessageListener
}

func NewEmailListener(listener dependencies.IMessageListener, emailSender usecases.SendEmailUseCase) EmailListener {
	return EmailListener{
		listener:    listener,
		emailSender: emailSender,
	}
}

func (listener EmailListener) ReceiveMessages(stopSig chan os.Signal) {
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
				message, err := messagecontracts.SendEmailMessage{}.FromBinary(msg.GetContent())

				if err != nil {
					msg.Nack(false)
					continue
				}

				email := models.Email{
					Receiver: message.Receiver,
					Subject:  message.Subject,
					Body:     message.Body,
				}

				err = listener.emailSender.Execute(email)

				if err != nil {
					msg.Nack(true)
					continue
				}

				msg.Ack()
			}
		}
	}()
}
