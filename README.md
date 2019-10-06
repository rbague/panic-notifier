# PANIC-NOTIFIER

[![GoDoc](https://godoc.org/github.com/rbague/panic-notifier/integration?status.svg)](https://godoc.org/github.com/rbague/panic-notifier/integration)
[![Go Report Card](https://goreportcard.com/badge/github.com/rbague/panic-notifier)](https://goreportcard.com/report/github.com/rbague/panic-notifier)

Panic notifier is based in the [Exception Notification](https://github.com/smartinez87/exception_notification/) ruby gem. You should definitely give it a look.

---

The purpose of this library is to receive a notification into an integration everytime a webserver panics.
Right now, the built-in notifiers can deliver notifications to Slack, email and via custom WebHooks.

## Getting Started

```go
go get github.com/rbague/panic-notifier
```

### Usage

To start using it with [go-chi/chi](https://github.com/go-chi/chi), [gorilla/mux](https://github.com/gorilla/mux) or similar:

```go
package main

import "github.com/rbague/panic-notifier"

func main() {
    r := chi.NewRouter() // github.com/go-chi/chi
    r := mux.NewRouter() // github.com/gorilla/mux

    s := integration.NewSlack("INCOMING_WEBHOOK_URL")
    r.Use(notifier.Middleware(s))

    ...
}
```

## Integrations

panic-notifier relies on integrations to deliver the notifications to the different services, here are the default ones:

-   [Slack](#slack)
-   [Email](#email)
-   [WebHooks](#webhooks)

But you could easily implement your [custom integration](#custom-integration).

### Slack

This integration delivers the notification to a slack channel using the [slack-webhook](https://github.com/rbague/slack-webhook) library.

#### Usage

To start using it you only need to provide the webhook url.

```go
package main

import "github.com/rbague/panic-notifier"

func main() {
    r := chi.NewRouter() // github.com/go-chi/chi

    slack := integration.NewSlack("INCOMING_WEBHOOK_URL")
    r.Use(notifier.Middleware(slack))

    ...
}
```

For production uses, recommend calling the following method with a `*http.Client` so it does not use the default client to deliver the notification.

```go
slack.SetHTTPClient(*http.Client)
```

### Email

This integration delivers the notification via email using SMTP.

#### Usage

To start using it with the default configuration you only need to provide a GMail email and password:

```go
package main

import "github.com/rbague/panic-notifier"

func main() {
    r := chi.NewRouter() // github.com/go-chi/chi

    email := integration.NewDefaultEmail("YOUR_EMAIL@gmail.com", "YOUR_PASSWORD")
    r.Use(notifier.Middleware(email))

    ...
}
```

The default configuration uses gmail's SMTP server (smtp.gmail.com) with SSL enabled (port: 465).
But this canbe easily overridden by providing your on configuration:

```go
package main

import "github.com/rbague/panic-notifier"

func main() {
    r := chi.NewRouter() // github.com/go-chi/chi

    config := &SMTPConfig{
		Addr:     "smtp.example.org",
		Port:     587,
		User:     "YOUR_EMAIL@gmail.com",
		Password: "YOUR_PASSWORD",
	}
    email := integration.NewEmail(config)
    r.Use(notifier.Middleware(email))

    ...
}
```

### WebHooks

This integration delivers the notification via the HTTP protocol.
Right now each request is sent using the `POST` method.

#### Usage

To start using it you only need to provide the webhook url.

```go
package main

import "github.com/rbague/panic-notifier"

func main() {
    r := chi.NewRouter() // github.com/go-chi/chi

    wh := integration.NewWebHook("WEBHOOK_URL")
    r.Use(notifier.Middleware(wh))

    ...
}
```

The WebHook type exposes its `Client` field, so the http.Client can be changed for production uses.

### Custom integration

Creating a custom integration is as easy as implementing the integration.Integration interface.

```go
// Integration is the interface used to deliver a notification of a panic
type Integration interface {
	// StackTraceLines return the number of stack trace lines
	// to be sent in each notification
	StackTraceLines() int

	// Deliver delivers the given notification and returns an error if any.
	Deliver(*Notification) error
}
```

## TODO

Add integrations for:

-   [x] Slack
-   [x] Email
-   [x] WebHooks

Other:

-   [ ] Add option to send custom data
-   [ ] Add tests
-   [ ] Allow custom webhooks to use other methods than POST

## License

Copyright (c) 2018 Roger Bagué Martí, released under the [MIT license](http://www.opensource.org/licenses/MIT).
