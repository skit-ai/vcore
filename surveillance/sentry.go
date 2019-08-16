package surveillance

import (
	"github.com/getsentry/sentry-go"
	"github.com/getsentry/sentry-go/http"
	"net/http"
	"os"
	"github.com/Vernacular-ai/vcore/errors"
	"github.com/Vernacular-ai/vcore/log"
)

type Sentry struct {
	client  *sentry.Client
	handler *sentryhttp.Handler
}

func initSentry() (client *Sentry) {
	dsn := os.Getenv("SENTRY_DSN")
	if dsn != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn: dsn,
			// Use async transport. Which is set by default. Use Sync transport for testing.
			//Transport: sentry.NewHTTPSyncTransport(),

			// Enable debugging to check connectivity
			//Debug: true,
		}); err != nil {
			log.Warnf("Could not initialize sentry with DSN: %s", dsn)
			client = &Sentry{nil, nil}
		} else {
			client = &Sentry{
				sentry.CurrentHub().Client(),
				sentryhttp.New(sentryhttp.Options{Repanic: true}),
			}
		}
	} else {
		log.Warnf("Could not initialize sentry with DSN: %s", dsn)
		client = &Sentry{nil, nil}
	}
	return
}

//func initSentry() *Sentry {
//	dsn := os.Getenv("SENTRY_DSN")
//	if dsn != "" {
//		if client, err := raven.New(dsn); err != nil {
//			return &Sentry{nil}
//		} else {
//			_client := &Sentry{client}
//			_client.client.SetRelease(os.Getenv("SENTRY_RELEASE"))
//			_client.client.SetEnvironment(os.Getenv("SENTRY_ENVIRONMENT"))
//			return _client
//		}
//	} else {
//		return &Sentry{nil}
//	}
//}

var (
	SentryClient = initSentry()
)

// Handles an error by capturing it on Sentry and logging the same on STDOUT
func (wrapper *Sentry) Capture(err error, _panic bool) {
	if err != nil {
		if wrapper.client != nil {
			// Capture error asynchronously
			sentry.WithScope(func(scope *sentry.Scope) {

				// Setting the stacktrace of the error as an extra
				scope.SetExtras(map[string]interface{}{
					"stacktrace": errors.Stacktrace(err),
				})

				// Determining the tags(if any) set on the error
				scope.SetTags(errors.Tags(err))

				// Capturing the error on Sentry
				eventId := sentry.CaptureException(err)
				log.Errorf(err, "Error captured in sentry with the event ID `%s`", *eventId)
			})
		} else {
			// Log the error sans sentry's event ID information
			log.Error(err)
		}

		if _panic {
			panic(err)
		}
	}
}

// Wrapper over sentry-go/http#HandleFunc
// Only calls the sentry handler if sentry was successfully initialized
func (wrapper *Sentry) HandleFunc(handler http.HandlerFunc) http.HandlerFunc {
	if wrapper.handler != nil {
		// If the sentry handler was initialized, call it's HandleFunc function
		return wrapper.handler.HandleFunc(handler)
	} else {
		// Simply return the handler in case the sentry handler was not initialized
		return handler
	}
}
