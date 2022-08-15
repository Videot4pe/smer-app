package mailer

import (
	"backend/pkg/logging"
	"gopkg.in/gomail.v2"
	"sync"
)

type Mailer struct {
	Sender SenderConfig
	Dialer *gomail.Dialer
	Logger *logging.Logger
}

type SenderConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type Mail struct {
	Username    string
	Subject     string
	Text        string
	Attachments []string
}

var instance *Mailer
var once sync.Once

func GetMailer(sender SenderConfig, logger *logging.Logger) *Mailer {
	once.Do(func() {
		dialer := gomail.NewDialer(sender.Host, sender.Port, sender.Username, sender.Password)
		instance = &Mailer{
			Dialer: dialer,
			Sender: sender,
			Logger: logger,
		}
	})
	return instance
}

func (m *Mailer) Send(mail Mail) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.Sender.Username)
	msg.SetHeader("To", mail.Username)
	msg.SetHeader("Subject", mail.Subject)
	msg.SetBody("text/html", mail.Text)

	for _, attachment := range mail.Attachments {
		msg.Attach(attachment)
	}

	if err := m.Dialer.DialAndSend(msg); err != nil {
		m.Logger.Error(err)
		return err
	}
	return nil
}
