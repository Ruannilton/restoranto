package presentation

import (
	"os"

	"github.com/Ruannilton/contact-verification-service/domain/dependencies"
	"github.com/Ruannilton/contact-verification-service/domain/usecases"
	messagecontracts "github.com/Ruannilton/go-msg-contracts/pkg/message_contracts"
	"github.com/Ruannilton/go-msg-contracts/pkg/queues"
)

const VerifyEmailListenerName = "contact-verification-verify-email-listener"

type VerifyEmailListener struct {
	listener dependencies.IMessageListener
	useCase  usecases.VerifyEmailUseCase
}

func NewVerifyEmailListener(listener dependencies.IMessageListener, useCase usecases.VerifyEmailUseCase) VerifyEmailListener {
	return VerifyEmailListener{
		listener: listener,
		useCase:  useCase,
	}
}

func (emailListener VerifyEmailListener) ReceiveMessages(stopSig chan os.Signal) {

	defer emailListener.listener.Close()

	messageChannel, err := emailListener.listener.Listen(queues.EmailValidationQueue)

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
				message, err := messagecontracts.VerifyEmailMessage{}.FromBinary(msg.GetContent())

				if err != nil {
					//TODO: handle error
					msg.Nack(false)
					continue
				}

				err = emailListener.useCase.Execute(message)

				if err != nil {
					//TODO: handle error
					msg.Nack(true)
					continue
				}

				msg.Ack()
			}
		}
	}()
}
