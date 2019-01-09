package integration

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// WebHook implements the Integration interface
// and delivers the notifications over the HTTP protocol
type WebHook struct {
	URL    string
	Client *http.Client
}

// NewWebHook returns a webhook integration
func NewWebHook(url string) *WebHook {
	return &WebHook{
		URL:    url,
		Client: http.DefaultClient,
	}
}

func (WebHook) StackTraceLines() int { return 15 }

func (w WebHook) Deliver(n *Notification) error {
	short := *n
	short.Stack = short.Stack[:w.StackTraceLines()]

	resp, err := w.Client.PostForm(w.URL, short.toURLValues())
	if err != nil {
		return fmt.Errorf("could not execute http request: %v", err)
	}
	defer resp.Body.Close()
	io.Copy(ioutil.Discard, resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("http request not successful got: %v", resp.Status)
	}

	return nil
}
