package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/sangtandoan/subscription_tracker/internal/config"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/apperror"
	"gopkg.in/gomail.v2"
)

type SMTPMailder struct {
	config *config.MailerConfig
	dialer *gomail.Dialer
}

func NewSMTPMailer(config *config.MailerConfig) *SMTPMailder {
	dialer := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)

	return &SMTPMailder{
		config,
		dialer,
	}
}

func (m *SMTPMailder) SendWithRetry(req *SendRequest, retries int) error {
	message, err := m.processMessage(req)
	if err != nil {
		return err
	}

	for i := range retries {
		if err := m.dialer.DialAndSend(message); err != nil {
			fmt.Println(err)

			if i != retries-1 {
				time.Sleep(time.Second * time.Duration((i + 1))) // exponential backkoff
			}
			continue
		}
		fmt.Println("send email ok")
		return nil
	}

	return apperror.ErrSendEmail
}

func (m *SMTPMailder) Send(req *SendRequest) error {
	message, err := m.processMessage(req)
	if err != nil {
		return err
	}

	if err := m.dialer.DialAndSend(message); err != nil {
		return err
	}

	return nil
}

func (m *SMTPMailder) processMessage(req *SendRequest) (*gomail.Message, error) {
	// parse and check if data is right with template option
	data, err := getData(req.Template, req.Data)
	if err != nil {
		return nil, err
	}

	// get template for template option
	temp := getMailTemplate(req.Template)

	// get subject and body base on template to send email
	err = parseTemplate(temp, data)
	if err != nil {
		return nil, err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", fmt.Sprintf("Subdub <%s>", m.config.From))
	message.SetHeader("To", req.To...)
	message.SetHeader("Subject", temp.Subject)
	message.SetBody("text/html", temp.Body)

	return message, nil
}

func getData(opt MailTemplateOption, data any) (any, error) {
	switch opt {
	case RemindTemplate:
		if data, ok := data.(RemindData); ok {
			data.Name = strings.ReplaceAll(data.Name, "subscription", "")
			data.Name = strings.ToUpper(data.Name)
			return data, nil
		}
		return nil, apperror.ErrInvalidEmailData
	}

	return nil, nil
}

func getMailTemplate(opt MailTemplateOption) *MailTemplate {
	var temp MailTemplate
	switch opt {
	case RemindTemplate:
		temp.Path = "remind-email.tmpl"
	}

	return &temp
}

//go:embed "templates/*"
var FS embed.FS

func parseTemplate(temp *MailTemplate, data any) error {
	t, err := template.ParseFS(FS, "templates/"+temp.Path)
	if err != nil {
		return err
	}

	var subject bytes.Buffer
	if err := t.ExecuteTemplate(&subject, "subject", data); err != nil {
		return err
	}

	var body bytes.Buffer
	if err := t.ExecuteTemplate(&body, "body", data); err != nil {
		return err
	}

	temp.Subject = subject.String()
	temp.Body = body.String()

	return nil
}
