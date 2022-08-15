package auth

import (
	"backend/internal/config"
	"backend/pkg/logging"
	"backend/pkg/mailer"
	"bytes"
	"html/template"
	"os"
)

type MailerAuth struct {
	Client *mailer.Mailer
	Logger *logging.Logger
}

const (
	EmailConfirmationTemplate = "/templates/email-confirmation.html"
)

type EmailConfirmationParams struct {
	Name  string
	Email string
	Link  string
}

func GetMailerAuth(cfg *config.Config, logger *logging.Logger) *MailerAuth {
	sender := mailer.SenderConfig{
		Host:     cfg.Mailer.Host,
		Port:     cfg.Mailer.Port,
		Username: cfg.Mailer.Username,
		Password: cfg.Mailer.Password,
	}

	mailClient := mailer.GetMailer(sender, logger)

	return &MailerAuth{
		Client: mailClient,
		Logger: logger,
	}
}

func (ma *MailerAuth) SendMail(username string, subject string, tmp string, params interface{}) error {
	templateString, err := ma.GetTemplate(tmp, params)
	if err != nil {
		return err
	}

	mail := mailer.Mail{
		Username: username,
		Subject:  subject,
		Text:     templateString,
	}

	err = ma.Client.Send(mail)
	if err != nil {
		return err
	}
	return nil
}

func (ma *MailerAuth) GetTemplate(tmp string, params interface{}) (string, error) {
	wd, err := os.Getwd()
	// TODO fix path
	t, err := template.ParseFiles(wd + tmp)
	if err != nil {
		ma.Logger.Error(err)
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, params); err != nil {
		ma.Logger.Error(err)
		return "", err
	}

	result := tpl.String()
	if err != nil {
		ma.Logger.Error(err)
		return "", err
	}

	return result, nil
}
