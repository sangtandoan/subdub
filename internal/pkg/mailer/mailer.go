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

type Mailer interface {
	Send(req *SendRequest) error
	SendWithRetry(req *SendRequest, retries int) error
}

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

type MailTemplateOption int

const (
	RemindTemplate MailTemplateOption = iota
)

type RemindData struct {
	RenewalDate time.Time
	Name        string
	Email       string
	NumDays     int
}

type SendRequest struct {
	To       []string
	Data     any
	Template MailTemplateOption
}

type MailTemplate struct {
	Subject string
	Body    string
	Path    string
}

func (m *SMTPMailder) SendWithRetry(req *SendRequest, retries int) error {
	for i := range retries {
		err := m.Send(req)
		if err != nil {
			fmt.Printf("Failed to send email %d/%d times\n", i+1, retries)
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
	// parse and check if data is right with template option
	data, err := getData(req.Template, req.Data)
	if err != nil {
		return err
	}

	// get template for template option
	temp := getMailTemplate(req.Template)

	// get subject and body base on template to send email
	err = parseTemplate(temp, data)
	if err != nil {
		return err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", m.config.From)
	message.SetHeader("To", req.To...)
	message.SetHeader("Subject", temp.Subject)
	message.SetBody("text/html", temp.Body)

	if err := m.dialer.DialAndSend(message); err != nil {
		return err
	}

	return nil
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
