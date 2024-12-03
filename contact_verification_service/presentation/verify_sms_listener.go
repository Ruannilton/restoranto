package presentation

import (
	"os"

	"github.com/Ruannilton/contact-verification-service/domain/dependencies"
	"github.com/Ruannilton/contact-verification-service/domain/usecases"
	messagecontracts "github.com/Ruannilton/go-msg-contracts/pkg/message_contracts"
	"github.com/Ruannilton/go-msg-contracts/pkg/queues"
)

const VerifySMSListenerName = "contact-verification-verify-sms-listener"

type VerifySMSListener struct {
	listener dependencies.IMessageListener
	useCase  usecases.VerifySMSUseCase
}

func NewVerifySMSListener(listener dependencies.IMessageListener, useCase usecases.VerifySMSUseCase) VerifySMSListener {
	return VerifySMSListener{
		listener: listener,
		useCase:  useCase,
	}
}

func (smsListener VerifySMSListener) ReceiveMessages(stopSig chan os.Signal) {
	defer smsListener.listener.Close()

	messageChannel, err := smsListener.listener.Listen(queues.SMSValidationQueue)

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
				message, err := messagecontracts.VerifySMSMessage{}.FromBinary(msg.GetContent())

				if err != nil {
					//TODO: handle error
					_ = msg.Nack(false)
					continue
				}

				err = smsListener.useCase.Execute(message)

				if err != nil {
					//TODO: handle error
					_ = msg.Nack(true)
					continue
				}

				_ = msg.Ack()
			}
		}
	}()
}
