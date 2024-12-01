package presentation

import (
	"os"

	usecase "github.com/Ruannilton/event-handler-service/domain/use_case"
	"github.com/Ruannilton/go-msg-contracts/pkg/events"
	"github.com/Ruannilton/go-msg-contracts/pkg/queues"
	amqp "github.com/rabbitmq/amqp091-go"
)

const CostumerCreatedListenerName = "event-handler-costumer-created-listener"

type CostumerCreatedListener struct {
	mqConnection *amqp.Connection
	useCase      usecase.CostumerCreatedUseCase
}

func NewCostumerCreatedListener(mqConnection *amqp.Connection, useCase usecase.CostumerCreatedUseCase) CostumerCreatedListener {
	return CostumerCreatedListener{
		mqConnection: mqConnection,
		useCase:      useCase,
	}
}

func (listener CostumerCreatedListener) ReceiveMessages(stopSig chan os.Signal) {
	channel, err := listener.mqConnection.Channel()

	if err != nil {
		//TODO: handle error
		return
	}

	defer channel.Close()

	messageChannel, err := channel.Consume(queues.CostumerCreatedEventQueue, CostumerCreatedListenerName, false, false, false, false, nil)

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
				message, err := desserializeMessage[events.CostumerCreatedEvent](msg.Body)

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
