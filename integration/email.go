package integration

import (
	"bytes"
	htmlTemplate "html/template"
	"io"
	textTemplate "text/template"

	"gopkg.in/gomail.v2"
)

// Email implements the Integration interface and uses
// SMTP to deliver the notifications
type Email struct {
	d       *gomail.Dialer
	headers map[string][]string

	text *textTemplate.Template
	html *htmlTemplate.Template
}

// SMTPConfig contains the SMTP configuration
type SMTPConfig struct {
	Addr     string
	Port     int
	User     string
	Password string
}

// NewDefaultEmail returns an email integration with the default settings.
// The default settings use gmail's smtp and ssl
func NewDefaultEmail(email, password string) *Email {
	return NewEmail(&SMTPConfig{
		Addr:     "smtp.gmail.com",
		Port:     465,
		User:     email,
		Password: password,
	})
}

// NewEmail returns an email integration
func NewEmail(cfg *SMTPConfig) *Email {
	d := gomail.NewDialer(cfg.Addr, cfg.Port, cfg.User, cfg.Password)
	h := map[string][]string{
		"From":    {cfg.User},
		"To":      {cfg.User},
		"Subject": {"Server panic"},
	}

	tt := textTemplate.Must(textTemplate.ParseFiles("integration/templates/email.txt"))
	ht := htmlTemplate.Must(htmlTemplate.ParseFiles("integration/templates/email.html"))

	return &Email{d: d, headers: h, text: tt, html: ht}
}

func (Email) StackTraceLines() int {
	return 15
}

func (e Email) Deliver(n *Notification) error {
	short := *n
	short.Stack = short.Stack[:e.StackTraceLines()]

	m := gomail.NewMessage()
	m.SetHeaders(e.headers)

	var buf bytes.Buffer
	e.text.Execute(&buf, short)
	m.SetBody("text/plain", buf.String())

	m.AddAlternativeWriter("text/html", func(w io.Writer) error {
		return e.html.Execute(w, short)
	})

	return e.d.DialAndSend(m)
}
