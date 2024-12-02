package main

import (
	"fmt"
	"os"

	"github.com/Ruannilton/event-handler-service/domain/dependencies"
	usecase "github.com/Ruannilton/event-handler-service/domain/use_case"
	"github.com/Ruannilton/event-handler-service/infrastructure"
	"github.com/Ruannilton/event-handler-service/presentation"
	amqp "github.com/rabbitmq/amqp091-go"
)

type EventHandlerApp struct {
	stopSig                     chan os.Signal
	rmqpConnection              *amqp.Connection
	messagePublisher            dependencies.IMessagePublisher
	costumerCreatedUseCase      usecase.CostumerCreatedUseCase
	costumerCreatedPresentation presentation.CostumerCreatedListener
}

func NewEventHandlerApp(stopSig chan os.Signal) (EventHandlerApp, error) {
	app := EventHandlerApp{
		stopSig: stopSig,
	}

	err := app.initInfrastructure()

	if err != nil {
		return EventHandlerApp{}, err
	}

	app.initUseCases()
	app.initPresentation()

	return app, nil
}

func (app *EventHandlerApp) initInfrastructure() error {
	rabbitConnectionString := os.Getenv("RABBITMQ_URL")
	rmqConn, err := amqp.Dial(rabbitConnectionString)
	if err != nil {
		fmt.Println("Failed to connect with rabbitmq")
		return err
	}

	publisher := infrastructure.NewMessagePublisher(rmqConn)
	app.messagePublisher = publisher
	app.rmqpConnection = rmqConn
	return nil
}

func (app *EventHandlerApp) initUseCases() {
	costumerCreated := usecase.NewCostumerCreatedUseCase(app.messagePublisher)

	app.costumerCreatedUseCase = costumerCreated
}

func (app *EventHandlerApp) initPresentation() {
	costumerCreated := presentation.NewCostumerCreatedListener(app.rmqpConnection, app.costumerCreatedUseCase)

	app.costumerCreatedPresentation = costumerCreated
}

func (app *EventHandlerApp) StartWorkers() {
	go func() {
		app.costumerCreatedPresentation.ReceiveMessages(app.stopSig)
	}()
}

func (app *EventHandlerApp) Close() {
	app.rmqpConnection.Close()
}
