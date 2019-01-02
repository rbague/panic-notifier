package integration

import (
	"net/url"
	"strings"
)

// Integration is the interface used to deliver a notification of a panic
type Integration interface {
	// StackTraceLines return the number of stack trace lines
	// to be displayed in each notification
	StackTraceLines() int

	// Deliver delivers the given notification into its integration.
	Deliver(*Notification) error
}

// Notification contains information about what caused the panic
type Notification struct {
	Err     error
	Host    string // the machine's hostname
	Stack   []string
	Request *Request
}

// Request holds information about the request that made the panic
type Request struct {
	Method string
	URI    string // the relative url
}

func (n Notification) toURLValues() url.Values {
	values := url.Values{}
	values.Set("error", n.Err.Error())
	values.Set("hostname", n.Host)
	values.Set("stacktrace", strings.Join(n.Stack, "\n"))
	values.Set("request_method", n.Request.Method)
	values.Set("request_uri", n.Request.URI)
	return values
}
