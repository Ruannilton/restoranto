package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Ruannilton/contact-verification-service/domain/dependencies"
	"github.com/Ruannilton/contact-verification-service/domain/usecases"
	"github.com/Ruannilton/contact-verification-service/infrastructure"
	"github.com/Ruannilton/contact-verification-service/presentation"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ContactVerificationApp struct {
	stopSig             chan os.Signal
	rmqpConnection      *amqp.Connection
	messagePublisher    dependencies.IMessagePublisher
	distributedCache    dependencies.IDistributedCacheRepository
	verifyEmailUseCase  usecases.VerifyEmailUseCase
	verifySmsUseCase    usecases.VerifySMSUseCase
	verifyEmailListener presentation.VerifyEmailListener
	verifySmsListener   presentation.VerifySMSListener
}

func NewContactVerificationApp(stopSig chan os.Signal) (ContactVerificationApp, error) {
	app := ContactVerificationApp{stopSig: stopSig}
	infraErr := app.initInfrastructure()

	if infraErr != nil {
		return ContactVerificationApp{}, infraErr
	}

	app.initUseCases()
	app.initPresentation()
	return app, nil
}

func (app *ContactVerificationApp) initInfrastructure() error {
	rabbitConnectionString := os.Getenv("RABBITMQ_URL")
	rmqConn, err := amqp.Dial(rabbitConnectionString)
	if err != nil {
		fmt.Println("Failed to connect with rabbitmq")
		return err
	}
	publisher := infrastructure.NewMessagePublisher(rmqConn)
	cache := infrastructure.NewRedisCache(context.Background())

	app.messagePublisher = publisher
	app.rmqpConnection = rmqConn
	app.distributedCache = cache
	return nil
}

func (app *ContactVerificationApp) initUseCases() {
	emailUseCase := usecases.NewVerifyEmailUseCase(app.messagePublisher, app.distributedCache)
	smsUseCase := usecases.NewVerifySMSUseCase(app.messagePublisher, app.distributedCache)

	app.verifyEmailUseCase = emailUseCase
	app.verifySmsUseCase = smsUseCase
}

func (app *ContactVerificationApp) initPresentation() {
	emailPresentation := presentation.NewVerifyEmailListener(app.rmqpConnection, app.verifyEmailUseCase)
	smsPresentation := presentation.NewVerifySMSListener(app.rmqpConnection, app.verifySmsUseCase)

	app.verifyEmailListener = emailPresentation
	app.verifySmsListener = smsPresentation
}

func (app *ContactVerificationApp) StartWorkers() {
	go func() {
		app.verifyEmailListener.ReceiveMessages(app.stopSig)
	}()
	go func() {
		app.verifySmsListener.ReceiveMessages(app.stopSig)
	}()
}

func (app *ContactVerificationApp) Close() {
	app.rmqpConnection.Close()
}
