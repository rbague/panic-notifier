package notifier

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/rbague/panic-notifier/integration"
)

// Middleware used to deliver a notification to every provided integration
// every time the webserver panics
func Middleware(ii ...integration.Integration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if len(ii) > 0 {
				defer func() {
					if err := recover(); err != nil {
						defer panic(err)

						h, _ := os.Hostname()
						n := &integration.Notification{
							Err:   asError(err),
							Host:  h,
							Stack: getStackTrace(),
							Request: &integration.Request{
								Method: r.Method,
								URI:    r.RequestURI,
							},
						}

						var err error
						for _, i := range ii {
							err = i.Deliver(n)
							if err != nil {
								log.Printf("could not deliver notification through %s: %v", reflect.TypeOf(i).String(), err)
							}
						}
					}
				}()
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func asError(err interface{}) error {
	if e, ok := err.(error); ok {
		return e
	}
	return fmt.Errorf("%v", err)
}

func getStackTrace() []string {
	s := debug.Stack()
	return strings.Split(string(s[:]), "\n")
}
