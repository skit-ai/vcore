package sentry

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type options struct {
	Repanic  bool
	ReportOn ReportOn
}

func BuildOptions(ff ...Option) options {
	opts := options{
		ReportOn: ReportAlways,
	}

	for _, f := range ff {
		f(&opts)
	}

	return opts
}

// Option configures reporting behavior.
type Option func(*options)

// WithRepanic configures whether to panic again after recovering from
// a panic. Use this option if you have other panic handlers.
func WithRepanic(b bool) Option {
	return func(o *options) {
		o.Repanic = b
	}
}

// WithReportOn configures whether to report on errors.
func WithReportOn(r ReportOn) Option {
	return func(o *options) {
		o.ReportOn = r
	}
}

// ReportOn decides error should be reported to sentry.
type ReportOn func(error) bool

// isContextCanceledError checks if an error is context.Canceled or derived from it
func isContextCanceledError(err error) bool {
	if err == context.Canceled {
		return true
	}

	// Check if it's a gRPC Canceled error
	if st, ok := status.FromError(err); ok && st.Code() == codes.Canceled {
		return true
	}

	// Check if it wraps context.Canceled
	return errors.Is(err, context.Canceled)
}

// ReportAlways returns true if err is non-nil and not a context canceled error.
func ReportAlways(err error) bool {
	return err != nil && !isContextCanceledError(err)
}

// ReportOnCodes returns true if error code matches on of the given codes.
func ReportOnCodes(cc ...codes.Code) ReportOn {
	cm := make(map[codes.Code]bool)
	for _, c := range cc {
		cm[c] = true
	}
	return func(err error) bool {
		return cm[status.Code(err)]
	}
}

// WrappedServerStream is a thin wrapper around grpc.ServerStream that allows modifying context.
type WrappedServerStream struct {
	grpc.ServerStream
	// WrappedContext is the wrapper's own Context. You can assign it.
	WrappedContext context.Context
}

// Context returns the wrapper's WrappedContext, overwriting the nested grpc.ServerStream.Context()
func (w *WrappedServerStream) Context() context.Context {
	return w.WrappedContext
}

// WrapServerStream returns a ServerStream that has the ability to overwrite context.
func WrapServerStream(stream grpc.ServerStream) *WrappedServerStream {
	if existing, ok := stream.(*WrappedServerStream); ok {
		return existing
	}
	return &WrappedServerStream{ServerStream: stream, WrappedContext: stream.Context()}
}
