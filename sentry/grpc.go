package sentry

import (
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

// ReportAlways returns true if err is non-nil.
func ReportAlways(err error) bool {
	return err != nil
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

