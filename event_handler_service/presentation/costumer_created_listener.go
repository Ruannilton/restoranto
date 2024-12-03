package presentation

import (
	"os"

	"github.com/Ruannilton/event-handler-service/domain/dependencies"
	usecase "github.com/Ruannilton/event-handler-service/domain/use_case"
	"github.com/Ruannilton/go-msg-contracts/pkg/events"
	"github.com/Ruannilton/go-msg-contracts/pkg/queues"
)

const CostumerCreatedListenerName = "event-handler-costumer-created-listener"

type CostumerCreatedListener struct {
	listener dependencies.IMessageListener
	useCase  usecase.CostumerCreatedUseCase
}

func NewCostumerCreatedListener(listener dependencies.IMessageListener, useCase usecase.CostumerCreatedUseCase) CostumerCreatedListener {
	return CostumerCreatedListener{
		listener: listener,
		useCase:  useCase,
	}
}

func (listener CostumerCreatedListener) ReceiveMessages(stopSig chan os.Signal) {
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
				message, err := desserializeMessage[events.CostumerCreatedEvent](msg.GetContent())

				if err != nil {
					//TODO: handle error
					msg.Nack(false)
					continue
				}

				err = listener.useCase.Execute(message)

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
