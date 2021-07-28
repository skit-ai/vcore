package sentry

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/julienschmidt/httprouter"
)


type Handler struct {
	repanic         bool
	waitForDelivery bool
	timeout         time.Duration
}

// HandleFunc wraps http.HandleFunc and recovers from caught panics.
func (h *Handler) HandleFunc(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		hub := sentry.GetHubFromContext(ctx)
		if hub == nil {
			hub = sentry.CurrentHub().Clone()
			ctx = sentry.SetHubOnContext(ctx, hub)
		}
		span := sentry.StartSpan(ctx, "http.server",
			sentry.TransactionName(fmt.Sprintf("%s %s", r.Method, r.URL.Path)),
			sentry.ContinueFromRequest(r),
		)
		defer span.Finish()
		// TODO(tracing): if the next handler.ServeHTTP panics, store
		// information on the transaction accordingly (status, tag,
		// level?, ...).
		r = r.WithContext(span.Context())
		hub.Scope().SetRequest(r)
		defer h.recoverWithSentry(hub, r)

		// TODO(tracing): use custom response writer to intercept
		// response. Use HTTP status to add tag to transaction; set span
		// status.
		handler.ServeHTTP(rw, r)
	}
}

// New returns a struct that provides Handle and HandleFunc methods
// that satisfy http.Handler and http.HandlerFunc interfaces.
func New(options sentryhttp.Options) *Handler {
	handler := Handler{
		repanic:         false,
		timeout:         time.Second * 2,
		waitForDelivery: false,
	}

	if options.Repanic {
		handler.repanic = true
	}

	if options.Timeout != 0 {
		handler.timeout = options.Timeout
	}

	if options.WaitForDelivery {
		handler.waitForDelivery = true
	}

	return &handler
}

// HandleFunc wraps http.HandleFunc and recovers from caught panics.
func (h *Handler) HandleHttpRouter(handler httprouter.Handle) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := r.Context()
		hub := sentry.GetHubFromContext(ctx)
		if hub == nil {
			hub = sentry.CurrentHub().Clone()
			ctx = sentry.SetHubOnContext(ctx, hub)
		}
		span := sentry.StartSpan(ctx, "http.server",
			sentry.TransactionName(fmt.Sprintf("%s %s", r.Method, r.URL.Path)),
			sentry.ContinueFromRequest(r),
		)
		defer span.Finish()
		// TODO(tracing): if the next handler.ServeHTTP panics, store
		// information on the transaction accordingly (status, tag,
		// level?, ...).
		r = r.WithContext(span.Context())
		hub.Scope().SetRequest(r)
		defer h.recoverWithSentry(hub, r)

		handler(rw, r, params)
	}
}

func (h *Handler) recoverWithSentry(hub *sentry.Hub, r *http.Request) {
	if err := recover(); err != nil {
		eventID := hub.RecoverWithContext(
			context.WithValue(r.Context(), sentry.RequestContextKey, r),
			err,
		)
		if eventID != nil && h.waitForDelivery {
			hub.Flush(h.timeout)
		}
		if h.repanic {
			panic(err)
		}
	}
}
