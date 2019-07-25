package surveillance

import (
	"github.com/getsentry/raven-go"
	"os"
	"vcore/errors"
)

type Sentry struct {
	client *raven.Client
}

func initSentry() *Sentry {
	dsn := os.Getenv("SENTRY_DSN")
	if dsn != "" {
		if client, err := raven.New(dsn); err != nil {
			return &Sentry{nil}
		} else {
			_client := &Sentry{client}
			_client.client.SetRelease(os.Getenv("SENTRY_RELEASE"))
			_client.client.SetEnvironment(os.Getenv("SENTRY_ENVIRONMENT"))
			return _client
		}
	} else {
		return &Sentry{nil}
	}
}

var (
	SentryClient = initSentry()
)

// Handles an error by capturing it on Sentry and logging the same on STDOUT
func (wrapper *Sentry) Capture(err error, _panic bool) {
	if err != nil {
		if _panic {
			panic(err)
		} else {
			errors.PrintStackTrace(err)
		}

		if wrapper.client != nil {
			// Determining the tags(if any) set on the error
			tags := errors.Tags(err)

			// Setting the stacktrace of the error as an extra
			extra := map[string]interface{}{
				"stacktrace": errors.Stacktrace(err),
			}
			_err := raven.WrapWithExtra(err, extra)

			// Not using CaptureError since it does not seem to be capturing the error on Sentry.
			wrapper.client.CaptureErrorAndWait(_err, tags)
		}
	}
}
