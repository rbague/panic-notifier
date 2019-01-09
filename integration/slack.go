package integration

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/rbague/slack-webhook"
)

// Slack implements the Integration interface and uses Slack's
// Incoming Webhooks API to deliver the notifications
type Slack struct {
	client *webhook.Client
}

// NewSlack returns a slack itegration
func NewSlack(url string) *Slack {
	return &Slack{client: webhook.NewClient(url)}
}

// SetHTTPClient sets the http.Client used to deliver the notification
func (s *Slack) SetHTTPClient(c *http.Client) {
	s.client.Client = c
}

func (Slack) StackTraceLines() int { return 15 }

var markdownReplacer = strings.NewReplacer("*", "", "`", "")

func (s Slack) Deliver(n *Notification) error {
	title := fmt.Sprintf("*Server panicked while* `%s` *to* `%s` *was being processed*", n.Request.Method, n.Request.URI)
	a := &webhook.Attachment{
		Fallback:   markdownReplacer.Replace(title),
		Text:       title,
		Color:      webhook.DangerColor,
		MarkdownIn: []string{"text", "fields"},
	}

	exception := webhook.Field{Title: "Error", Value: n.Err.Error()}
	hostname := webhook.Field{Title: "Hostname", Value: n.Host}

	stack := strings.Join(n.Stack[:s.StackTraceLines()], "\n")
	stacktrace := webhook.Field{
		Title: "Stack Trace",
		Value: fmt.Sprintf("```%s```", stack),
	}

	a.Fields = append(a.Fields, exception, hostname, stacktrace)
	p := &webhook.Payload{Attachments: []*webhook.Attachment{a}}
	return s.client.Send(p)
}
