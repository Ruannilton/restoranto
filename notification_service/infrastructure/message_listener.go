package infrastructure

import (
	"github.com/Ruannilton/notification-service/domain/dependencies"
	amqp "github.com/rabbitmq/amqp091-go"
)

const ListenerName = "contact-verification-listener"

type MessageInput struct {
	delivery amqp.Delivery
}

type MessageListener struct {
	mqConnection *amqp.Connection
	channel      *amqp.Channel
}

func NewMessageListener(mqConnection *amqp.Connection) (MessageListener, error) {
	channel, err := mqConnection.Channel()

	if err != nil {
		return MessageListener{}, err
	}

	return MessageListener{
		mqConnection: mqConnection,
		channel:      channel,
	}, nil
}

func (listener MessageListener) Listen(queueName string) (chan dependencies.IMessageInput, error) {

	queue, err := listener.channel.QueueDeclare(queueName, true, false, false, false, nil)

	if err != nil {
		// TODO: handle error
		return nil, err
	}

	messageChannel, err := listener.channel.Consume(queue.Name, ListenerName, false, false, false, false, nil)

	if err != nil {
		//TODO: handle error
		return nil, err
	}

	inputChannel := make(chan dependencies.IMessageInput)

	go func() {
		defer close(inputChannel)
		for msg := range messageChannel {
			inputChannel <- &MessageInput{delivery: msg}
		}
	}()

	return inputChannel, nil

}

func (listener MessageListener) Close() error {
	err := listener.channel.Close()

	if err != nil {
		//TODO: handle error
		return err
	}

	return nil
}

func (m MessageInput) GetContent() []byte {
	return m.delivery.Body
}

func (m MessageInput) Ack() error {
	return m.delivery.Ack(false)
}

func (m MessageInput) Nack(requeue bool) error {
	return m.delivery.Nack(false, requeue)
}
