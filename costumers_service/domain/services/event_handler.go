package services

import (
	"costumers-api/domain/dependencies"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Ruannilton/go-msg-contracts/pkg/events"
)

type EventHandler struct {
	eventRepository  dependencies.IEventRepository
	messagePublisher dependencies.IMessagePublisher
}

func NewEventHandler(
	eventRepository dependencies.IEventRepository,
	messagePublisher dependencies.IMessagePublisher,
) EventHandler {

	return EventHandler{
		eventRepository:  eventRepository,
		messagePublisher: messagePublisher,
	}
}

func (handler *EventHandler) RunJob(stopChannel chan os.Signal) {
	fmt.Println("Starting event handler worker")
	go func() {
		for {
			select {
			case <-stopChannel:
				return
			default:
				handler.run()
			}
		}
	}()
}

func (handler *EventHandler) run() {
	outboxMessages, err := handler.eventRepository.GetEvents(16)

	if err != nil {
		// TODO: handle error
		return
	}

	for _, message := range outboxMessages {
		var err error

		if message.Name == "UserCreatedEvent" {
			err = handler.handleCostumerCreatedEvent(string(message.Message))
			if err != nil {
				fmt.Println("Message not delivered:", err)
			}
		}

		if err == nil {
			handler.eventRepository.ProcessEvent(message)
		}

		// TODO: handle errors
	}
}

func (handler *EventHandler) handleCostumerCreatedEvent(body string) error {
	var event events.CostumerCreatedEvent

	if err := json.Unmarshal([]byte(body), &event); err != nil {
		return err
	}

	err := handler.messagePublisher.Publish(event)

	if err != nil {
		return err
	}

	return nil
}
