package infrastructure

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const defaultExchange = ""

type MessagePublisher struct {
	rabbitConn *amqp.Connection
}

func NewMessagePublisher(rabbitConn *amqp.Connection) MessagePublisher {
	return MessagePublisher{
		rabbitConn: rabbitConn,
	}
}

func (publisher MessagePublisher) Publish(message any, queueName string) error {

	channel, err := publisher.rabbitConn.Channel()
	if err != nil {
		// TODO: handle error
		return err
	}
	defer channel.Close()

	if err := channel.Confirm(false); err != nil {
		return fmt.Errorf("failed to enable publisher confirms: %v", err)
	}

	confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	queue, err := channel.QueueDeclare(queueName, true, false, false, false, nil)

	if err != nil {
		// TODO: handle error
		return err
	}

	bodyContent, err := json.Marshal(message)

	if err != nil {
		// TODO: handle error
		return err
	}

	err = channel.Publish(
		defaultExchange,
		queue.Name,
		true,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyContent,
		},
	)

	if err != nil {
		// TODO: handle error
		return err
	}

	confirm := <-confirms

	if !confirm.Ack {
		return fmt.Errorf("message not sent")
	}

	return nil
}
