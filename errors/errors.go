package errors

import (
	"fmt"
	_err "github.com/pkg/errors"
	"strings"
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
	msg   string
	cause error
	fatal bool
	tags  map[string]string
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

// Based on https://godoc.org/github.com/pkg/errors#hdr-Formatted_printing_of_errors
type stackTracer interface {
	StackTrace() _err.StackTrace
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
