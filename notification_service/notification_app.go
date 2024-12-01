package main

import (
	"fmt"
	"os"

	"github.com/Ruannilton/notification-service/domain/dependencies"
	usecases "github.com/Ruannilton/notification-service/domain/use_cases"
	"github.com/Ruannilton/notification-service/infrastructure"
	"github.com/Ruannilton/notification-service/presentation"
	amqp "github.com/rabbitmq/amqp091-go"
)

type NotificationApp struct {
	stopSig          chan os.Signal
	rmqpConnection   *amqp.Connection
	emailSender      dependencies.IEmailSender
	smsSender        dependencies.ISMSSender
	sendEmailUseCase usecases.SendEmailUseCase
	sendSMSUseCase   usecases.SendSMSUseCase
	emailListener    presentation.EmailListener
	smsListener      presentation.SMSListener
}

func NewNotificationApp(stopSig chan os.Signal) (NotificationApp, error) {
	app := NotificationApp{
		stopSig: stopSig,
	}

	err := app.initInfrastructure()

	if err != nil {
		return NotificationApp{}, err
	}

	app.initUseCases()
	app.initPresentation()

	return app, nil
}

func (app *NotificationApp) initInfrastructure() error {
	rabbitConnectionString := os.Getenv("RABBIT_MQ_CONNECTION")
	rmqConn, err := amqp.Dial(rabbitConnectionString)

	if err != nil {
		fmt.Println("Failed to connect with rabbitmq")
		return err
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	sender := os.Getenv("SMTP_SENDER")
	password := os.Getenv("SMTP_PASSWORD")

	emailSender := infrastructure.NewEmailSender(sender, password, smtpPort, smtpHost)
	smsSender := infrastructure.NewSmsSender()

	app.rmqpConnection = rmqConn
	app.emailSender = emailSender
	app.smsSender = smsSender

	return nil
}

func (app *NotificationApp) initUseCases() {
	sendEmail := usecases.NewSendEmailUseCase(app.emailSender)
	sendSms := usecases.NewSendSMSUseCase(app.smsSender)

	app.sendEmailUseCase = sendEmail
	app.sendSMSUseCase = sendSms
}

func (app *NotificationApp) initPresentation() {
	emailListener := presentation.NewEmailListener(app.rmqpConnection, app.sendEmailUseCase)
	smsListener := presentation.NewSMSListener(app.rmqpConnection, app.sendSMSUseCase)

	app.emailListener = emailListener
	app.smsListener = smsListener
}

func (app *NotificationApp) StartWorkers() {
	go func() {
		app.emailListener.ReceiveMessages(app.stopSig)
	}()

	go func() {
		app.smsListener.ReceiveMessages(app.stopSig)
	}()
}

func (app *NotificationApp) Close() {
	app.rmqpConnection.Close()
}
