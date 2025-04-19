package mailer

import (
	"time"
)

type Mailer interface {
	Send(req *SendRequest) error
	SendWithRetry(req *SendRequest, retries int) error
}

type SendRequest struct {
	Data     any
	To       []string
	Template MailTemplateOption
}

type MailTemplate struct {
	Subject string
	Body    string
	Path    string
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
