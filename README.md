# PANIC-NOTIFIER
Panic notifier is based in the [Exception Notification](https://github.com/smartinez87/exception_notification/) ruby gem. You should definitely give it a look.

---
The purpose of this library is to receive notification into an integration everytime a webserver panics.
Right now, the built-in notifiers can deliver notifications to Slack.

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
`notifier.Middleware` is a variadic funtion that receives Integration implementations as parameters.

## Integrations
Panic-notification uses Integrations to deliver the notification to different services. Currently, the only built'in integration is:
* [Slack](#slack)

But you can easily implement your own [custom integration](#custom-integration).

### Slack
This integrations send the notification using Slack's [Incoming Webhooks](https://api.slack.com/incoming-webhooks)

#### Usage


## TODO
Add integrations for:
- [x] Slack
- [ ] Email
- [ ] Telegram
- [ ] Custom webhooks

Other:
- [ ] Add option to send custom data
- [ ] Add tests

## License
Copyright (c) 2018 Roger Bagué Martí, released under the [MIT license](http://www.opensource.org/licenses/MIT).