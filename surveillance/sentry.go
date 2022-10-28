package surveillance

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/Vernacular-ai/vcore/errors"
	"github.com/Vernacular-ai/vcore/log"
	sentryWrapper "github.com/Vernacular-ai/vcore/sentry"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Sentry struct {
	client  *sentry.Client
	handler *sentryWrapper.Handler
}

func InitSentry(release string) (client *Sentry) {
	dsn := os.Getenv("SENTRY_DSN")

	sampleRate, _ := strconv.ParseFloat(os.Getenv("SENTRY_SAMPLING"), 64)

	if release == "" {
		release = os.Getenv("SENTRY_RELEASE")
	}
	if dsn != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              dsn,
			AttachStacktrace: true,
			// Use async transport. Which is set by default. Use Sync transport for testing.
			//Transport: sentry.NewHTTPSyncTransport(),

			// Enable debugging to check connectivity
			//Debug: true,
			Release:    release,
			SampleRate: sampleRate,

			Environment: os.Getenv("ENVIRONMENT"),
		}); err != nil {
			log.Warnf("Could not initialize sentry with DSN: %s", dsn)
			client = &Sentry{nil, nil}
		} else {
			client = &Sentry{
				sentry.CurrentHub().Client(),
				sentryWrapper.New(sentryhttp.Options{Repanic: true}),
			}
		}
	} else {
		log.Warnf("Could not initialize sentry with DSN: %s", dsn)
		client = &Sentry{nil, nil}
	}
	return
}

var (
	SentryClient = InitSentry("")
)

// Handles an error by capturing it on Sentry and logging the same on STDOUT
func (wrapper *Sentry) Capture(err error, _panic bool) sentry.EventID {
	eventID := new(sentry.EventID)
	if err != nil {
		// Do not log to sentry if the error is ignorable.
		// However, do log it to stdout
		if wrapper.client != nil && !errors.Ignore(err) {
			// Capture error asynchronously
			sentry.WithScope(func(scope *sentry.Scope) {

				// Setting the stacktrace of the error as an extra along with any other extras set in the error
				if extras := errors.Extras(err); extras != nil {
					scope.SetContext("extras", extras)

					// setExtras is deprecated
					// adding it for backward compatibility with vernacular's sentry
					scope.SetExtras(extras)
				}

				// Determining the tags(if any) set on the error
				scope.SetTags(errors.Tags(err))

				// Capturing the error on Sentry
				eventID = sentry.CaptureException(err)
				if eventID != nil {
					log.Errorf(err, "Error captured in sentry with the event ID `%s`", *eventID)
				}

				return
			})
		} else {
			// Log the error sans sentry's event ID information
			log.Error(err)
		}

		if _panic {
			panic(err)
		}
	}

	return *eventID
}

// Handles an error by capturing it on Sentry and logging the same on STDOUT
func (wrapper *Sentry) CaptureWithContext(c context.Context, err error, _panic bool) sentry.EventID {
	eventID := new(sentry.EventID)
	if err != nil {
		// Do not log to sentry if the error is ignorable.
		// However, do log it to stdout
		if wrapper.client != nil && !errors.Ignore(err) {
			// Capture error asynchronously
			if hub := sentry.GetHubFromContext(c); hub != nil {
				sentry.WithScope(func(scope *sentry.Scope) {
					// Setting the stacktrace of the error as an extra along with any other extras set in the error
					if extras := errors.Extras(err); extras != nil {
						scope.SetContext("extras", extras)

						// setExtras is deprecated
						// adding it for backward compatibility with vernacular's sentry
						scope.SetExtras(extras)
					}

					// Determining the tags(if any) set on the error
					scope.SetTags(errors.Tags(err))
				})

				// Capturing the error on Sentry
				eventID = hub.CaptureException(err)
				log.Errorf(err, "Error captured in sentry with the event ID `%s`", *eventID)
			} else {
				wrapper.Capture(err, _panic)
			}
		} else {
			// Log the error sans sentry's event ID information
			log.Error(err)
		}

		if _panic {
			panic(err)
		}
	}

	return *eventID
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

// Wrapper over sentry-go/http#HandleFunc
// Only calls the sentry handler if sentry was successfully initialized
func (wrapper *Sentry) HandleHttpRouter(handler httprouter.Handle) httprouter.Handle {
	if wrapper.handler != nil {
		// If the sentry handler was initialized, call it's HandleFunc function
		return wrapper.handler.HandleHttpRouter(handler)
	} else {
		// Simply return the handler in case the sentry handler was not initialized
		return handler
	}
}

// SentryMiddleware use directly with mux
// returns http.Handler to directly use with router
func (wrapper *Sentry) SentryMiddleware(next http.Handler) http.Handler {
	return wrapper.HandleFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// UnaryServerInterceptor is a grpc interceptor that reports errors and panics
// to sentry. It also sets *sentry.Hub to context.
func (wrapper *Sentry) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	opts := sentryWrapper.BuildOptions(sentryWrapper.WithRepanic(false))

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		hub := sentry.GetHubFromContext(ctx)
		if hub == nil {
			hub = sentry.CurrentHub().Clone()
			ctx = sentry.SetHubOnContext(ctx, hub)
		}

		defer func() {
			if r := recover(); r != nil {
				hub.RecoverWithContext(ctx, r)

				if opts.Repanic {
					panic(r)
				}

				err = status.Errorf(codes.Internal, "%s", r)
			}
		}()

		resp, err = handler(ctx, req)

		if opts.ReportOn(err) {
			hub.CaptureException(err)
		}

		return resp, err
	}
}

// StreamServerInterceptor returns a grpc interceptor that reports errors and panics
// to sentry. It also sets *sentry.Hub to context.
func (wrapper *Sentry) StreamServerInterceptor() grpc.StreamServerInterceptor {
	opts := sentryWrapper.BuildOptions(sentryWrapper.WithRepanic(false))

	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		ctx := stream.Context()
		hub := sentry.GetHubFromContext(ctx)
		if hub == nil {
			hub = sentry.CurrentHub().Clone()
			ctx = sentry.SetHubOnContext(ctx, hub)
		}

		defer func() {
			if r := recover(); r != nil {
				hub.RecoverWithContext(ctx, r)

				if opts.Repanic {
					panic(r)
				}

				_ = status.Errorf(codes.Internal, "%s", r)
			}
		}()

		wrapped := sentryWrapper.WrapServerStream(stream)
		wrapped.WrappedContext = ctx
		err := handler(srv, wrapped)

		if opts.ReportOn(err) {
			hub.CaptureException(err)
		}

		return err
	}
}
