package sentry

import (
	"context"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)


type Handler struct {
	repanic         bool
	waitForDelivery bool
	timeout         time.Duration
}

// HandleFunc wraps http.HandleFunc and recovers from caught panics.
func (h *Handler) HandleFunc(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		hub := sentry.CurrentHub().Clone()
		hub.Scope().SetRequest(sentry.Request{}.FromHTTPRequest(r))
		ctx := sentry.SetHubOnContext(
			r.Context(),
			hub,
		)
		defer h.recoverWithSentry(hub, r)
		handler(rw, r.WithContext(ctx))
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
		hub := sentry.CurrentHub().Clone()
		hub.Scope().SetRequest(sentry.Request{}.FromHTTPRequest(r))
		ctx := sentry.SetHubOnContext(
			r.Context(),
			hub,
		)
		defer h.recoverWithSentry(hub, r)
		handler(rw, r.WithContext(ctx), params)
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
