// The errors package is a wrapper around the brilliant
// https://github.com/pkg/errors library. The major difference between
// https://github.com/pkg/errors and this library is support for the following:
//
// • Wrapping of an error into a custom error so as to add a stacktrace.
//
// • Allows creation of errors with a stack which have no cause(nil cause).
//
// • Every error created supports the `fatality` interface which is meant to inform us if the error is fatal or not.
//
// • Every error created supports the `tagged` interface which returns any tags(`map[string]string`) associated with an error.
//
// • The stacktrace of the error will return the stactrace starting from the deepest cause.
//
// Create Error without cause
//
// This error is non-fatal. As marked by the last bool flag.
//
//    errors.NewError("Error without a cause", nil, false)
//
// Create Error with a cause
//
// Embedding a cause with the error.
//    cause := errors.NewError("Error without a cause", nil, false)
//    errWithCause := errors.NewError("Error with a cause", cause, false)
//
// Check if an error is fatal
//
// Sample below:
//     cause := errors.NewError("Error without a cause", nil, false)
//     if cause.Fatal{
//         fmt.Println("This error is fatal.")
//     }
//
// Get the stacktrace of the error and print it
//
// Sample below:
//    cause := errors.NewError("Error without a cause", nil, false)
//    errWithCause := errors.NewError("Error with a cause", cause, false)
//    fmt.Println(errWithCause.Stacktrace())
//
// Alternatively,
//
//    cause := errors.NewError("Error without a cause", nil, false)
//    errWithCause := errors.NewError("Error with a cause", cause, false)
//    errWithCause.PrintStackTrace()
package errors

import (
	"fmt"
	"net/http"
	"strings"

	_err "github.com/pkg/errors"
)

// Fatal is a condition to see if an error can be ignored or not.
// An error value has an Fatal condition if it implements the following
// interface:
//
//     type fatality interface {
//            Fatal() bool
//     }
//
// If the error does not implement Fatal, false will be returned.
// If the error is nil, false will be returned without further investigation.
// The logic will loop through the topmost error of the stack followed by all
// it's causes provided it implements the causer interface:
//
//	  type causer interface {
//			  Cause() error
//	  }
// If any one of the causes is fatal, the error is deemed fatal. i.e. irrecoverable
func Fatal(err error) (isFatal bool) {
	type fatality interface {
		Fatal() bool
	}

	// Keep going through all the errors in the stack until we hit one error which implements fatality
	// We use this first error to check if the error is fatal or not.
	for err != nil {
		if check, ok := err.(fatality); ok {
			isFatal = check.Fatal()
			break
		}

		// Going to the cause of the current error(if any)
		cause, ok := err.(causer)
		if !ok {
			// Since there is no cause of the current error, it is the root error(original error) that caused the issue
			// in the first place. Hence breaking the loop.
			break
		}

		err = cause.Cause()
	}

	return
}

// Custom error that implements:
// - cause interface from github.com/pkg/errors
// - error interface from go builtin
// - fatality interface from FSM
// It represents a rung in the chain of errors leading to the cause.
type rung struct {
	msg    string
	cause  error
	fatal  bool
	tags   map[string]string
	extras map[string]interface{}
	ignore bool
	code   int
}

func (e *rung) Error() (errorMsg string) {
	if e.msg != "" && e.cause != nil {
		errorMsg = fmt.Sprintf("%v \n\t==>> %v", e.msg, e.cause)
	} else if e.msg == "" && e.cause != nil {
		errorMsg = fmt.Sprintf("%v", e.cause)
	} else if e.msg != "" && e.cause == nil {
		errorMsg = fmt.Sprintf("%s", e.msg)
	}
	return
}

// Implementing the causer interface from github.com/pkg/errors
func (e *rung) Cause() error {
	return e.cause
}

func (e *rung) Fatal() bool {
	return e.fatal
}

func (e *rung) Tags() map[string]string {
	return e.tags
}

func (e *rung) Extras() map[string]interface{} {
	return e.extras
}

func (e *rung) Ignore() bool {
	return e.ignore
}

func (e *rung) Code() int {
	return e.code
}

// Creates an error which is chained with a cause
func NewError(_msg string, _cause error, _fatal bool) error {
	return NewErrorWithTags(_msg, _cause, _fatal, nil)
}

// Creates an error which is chained with a cause
func NewErrorWithTags(_msg string, _cause error, _fatal bool, _tags map[string]string) error {
	err := &rung{
		cause: _cause,
		msg:   _msg,
		fatal: _fatal,
		tags:  _tags,
	}
	return _err.WithStack(err)
}

// Creates an error which is chained with a cause
func NewErrorWithExtras(_msg string, _cause error, _fatal bool, _extras map[string]interface{}) error {
	err := &rung{
		cause:  _cause,
		msg:    _msg,
		fatal:  _fatal,
		tags:   nil,
		extras: _extras,
	}
	return _err.WithStack(err)
}

func NewErrorWithTagsAndExtras(_msg string, _cause error, _fatal bool, _tags map[string]string, _extras map[string]interface{}) error {
	err := &rung{
		cause:  _cause,
		msg:    _msg,
		fatal:  _fatal,
		tags:   _tags,
		extras: _extras,
	}
	return _err.WithStack(err)
}

// NewErrorToIgnore returns an error that informs loggers to ignore it
func NewErrorToIgnore(_msg string, _cause error) error {
	err := &rung{
		cause:  _cause,
		msg:    _msg,
		fatal:  false,
		ignore: true,
	}
	return _err.WithStack(err)
}

// NewErrorWithCode returns an error that has an int code associated with it
func NewErrorWithCode(_msg string, code int, _cause error) error {
	err := &rung{
		cause:  _cause,
		msg:    _msg,
		fatal:  false,
		ignore: false,
		code:   code,
	}
	return _err.WithStack(err)
}

// Based on https://godoc.org/github.com/pkg/errors#hdr-Formatted_printing_of_errors
type stackTracer interface {
	StackTrace() _err.StackTrace
}

// AddTagsToError checks if the input error implements a causer interface. If it does,
// it checks if the Cause is of internal type *rung. When this check also passes,
// it adds the input _tags to the existing tags information. In case the checks fail,
// the original error is returned.
func AddTagsToError(err error, _tags map[string]string) error {
	errCauser, ok := err.(causer)
	if !ok || errCauser == nil {
		return err
	}

	errAsRung, ok := errCauser.Cause().(*rung)
	if !ok || errAsRung == nil {
		return err
	}

	existingTags := errAsRung.Tags()
	if existingTags == nil {
		existingTags = make(map[string]string)
	}

	for key, value := range _tags {
		existingTags[key] = value
	}

	errAsRung.tags = existingTags
	return errAsRung
}

// AddExtrasToError checks if the input error implements a causer interface. If it does,
// it checks if the Cause is of internal type *rung. When this check also passes,
// it adds the input _extras to the existing extras information. In case the checks fail,
// the original error is returned.
func AddExtrasToError(err error, _extras map[string]interface{}) error {
	errCauser, ok := err.(causer)
	if !ok || errCauser == nil {
		return err
	}

	errAsRung, ok := errCauser.Cause().(*rung)
	if !ok || errAsRung == nil {
		return err
	}

	existingExtras := errAsRung.Extras()
	if existingExtras == nil {
		existingExtras = make(map[string]interface{})
	}

	for key, value := range _extras {
		existingExtras[key] = value
	}

	errAsRung.extras = existingExtras
	return errAsRung
}

// Determines the stacktrace of an error.
// It will retrieve the entire stacktrace starting from the original root cause
func Stacktrace(err error) string {
	// Printing the message of the original error
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%v\n", err))

	// Find the deepest element in the stack which implements the stackTracer interface
	var deepestStacktracer stackTracer
	for err != nil {
		if val, ok := err.(stackTracer); ok {
			deepestStacktracer = val
		}

		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}

	// Printing the entire stacktrace starting from the original cause of this issue
	if deepestStacktracer != nil {
		for _, f := range deepestStacktracer.StackTrace() {
			builder.WriteString(fmt.Sprintf("%+s:%d\n", f, f))
		}
	}
	return builder.String()
}

// Printing the stacktrace of an error.
// It will print the entire stacktrace starting from the original root cause
func PrintStackTrace(err error) {
	if err != nil {
		fmt.Println(Stacktrace(err))
	}
}

// Copying the causer interface from pkg/errors.
// This will be used to loop over the chain of causes leading up to the topmost error
type causer interface {
	Cause() error
}

func Tags(err error) (cumulativeTags map[string]string) {
	type tagged interface {
		Tags() map[string]string
	}

	// Keep going through all the errors in the stack and make a cumulative map of all the tags
	for err != nil {
		if check, ok := err.(tagged); ok {
			tagsSet := check.Tags()
			if tagsSet != nil {
				for k, v := range tagsSet {
					if cumulativeTags == nil {
						cumulativeTags = make(map[string]string)
					}
					// The highest error in the stack overrides the tag value set by the lower error in the stack
					if _, exists := cumulativeTags[k]; !exists {
						cumulativeTags[k] = v
					}
				}
			}
		}

		// Going to the cause of the current error(if any)
		cause, ok := err.(causer)
		if !ok {
			// Since there is no cause of the current error, it is the root error(original error) that caused the issue
			// in the first place. Hence breaking the loop.
			break
		}

		err = cause.Cause()
	}

	return
}

func Extras(err error) (cumulativeExtras map[string]interface{}) {
	type extra interface {
		Extras() map[string]interface{}
	}

	// Keep going through all the errors in the stack and make a cumulative map of all the tags
	for err != nil {
		if check, ok := err.(extra); ok {
			extrasSet := check.Extras()
			if extrasSet != nil {
				for k, v := range extrasSet {
					if cumulativeExtras == nil {
						cumulativeExtras = make(map[string]interface{})
					}
					// The highest error in the stack overrides the tag value set by the lower error in the stack
					if _, exists := cumulativeExtras[k]; !exists {
						cumulativeExtras[k] = v
					}
				}
			}
		}

		// Going to the cause of the current error(if any)
		cause, ok := err.(causer)
		if !ok {
			// Since there is no cause of the current error, it is the root error(original error) that caused the issue
			// in the first place. Hence breaking the loop.
			break
		}

		err = cause.Cause()
	}

	return
}

func Ignore(err error) bool {
	type ignore interface {
		Ignore() bool
	}

	// Keep going through all the errors in the stack and find if any error is supposed to be ignored
	for err != nil {
		if check, ok := err.(ignore); ok {
			ignore := check.Ignore()
			if ignore {
				return true
			}
		}

		// Going to the cause of the current error(if any)
		cause, ok := err.(causer)
		if !ok {
			// Since there is no cause of the current error, it is the root error(original error) that caused the issue
			// in the first place. Hence breaking the loop.
			break
		}

		err = cause.Cause()
	}

	return false
}

// Finds the deepest non-nil cause
func DeepestCause(err error) error {
	var cause causer
	var ok bool
	for err != nil {
		if cause, ok = err.(causer); !ok {
			// Since there is no cause of the current error, it is the root error(original error) that caused the issue
			// in the first place. Hence breaking the loop.
			return err
		}
		if cause.Cause() != nil {
			err = cause.Cause()
		} else {
			break
		}
	}
	return err
}

// This allows us to attach a response code with an error. This is achieved by implementing
// the errorCode interface:
//     type errorCode interface {
//            Code() bool
//     }
//
// If the error does not implement errorCode, 500 will be returned.
// If the error is nil, 200 will be returned without further investigation.
// The logic will loop through the topmost error of the stack followed by all
// it's causes provided it implements the causer interface:
//
//	  type causer interface {
//			  Cause() error
//	  }
// If any one of the causes implements errorCode and returns a non-zero Code(),
// that code is returned to the calling function.
func Code(err error, defaultCode int) (code int) {

	if err == nil {
		return http.StatusOK
	}

	type errorCode interface {
		Code() int
	}

	// Keep going through all the errors in the stack until we hit one error
	// which implements errorCode and has a non-zero error code.
	// We use this first error to return the error code.
	for err != nil {
		if check, ok := err.(errorCode); ok {
			if code = check.Code(); code != 0 {
				break
			}
		}

		// Going to the cause of the current error(if any)
		cause, ok := err.(causer)
		if !ok {
			// Since there is no cause of the current error, it is the root error(original error) that caused the issue
			// in the first place. Hence breaking the loop.
			break
		}

		err = cause.Cause()
	}

	if code <= 0 {
		if defaultCode <= 0 {
			code = http.StatusInternalServerError
		} else {
			code = defaultCode
		}
	}
	return
}
