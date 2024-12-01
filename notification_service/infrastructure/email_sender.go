package infrastructure

import (
	"net/smtp"

	"github.com/Ruannilton/notification-service/domain/models"
)

type EmailSender struct {
	auth     smtp.Auth
	smtpHost string
	smtpPort string
	sender   string
}

func NewEmailSender(sender string, password string, smtpPort string, smtpHost string) EmailSender {
	auth := smtp.PlainAuth("", sender, password, smtpHost)

	return EmailSender{
		auth:     auth,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		sender:   sender,
	}
}

func (client EmailSender) Send(e models.Email) error {

	message := []byte(e.Subject + "\n" + e.Body)

	err := smtp.SendMail(client.smtpHost+":"+client.smtpPort, client.auth, client.sender, []string{e.Receiver}, message)

	if err != nil {
		// TODO: handle err
		return err
	}

	return nil
}
