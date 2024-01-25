// Package utils implements helper functions
package utils

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"math"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/skit-ai/vcore/env"
	"github.com/skit-ai/vcore/errors"
	"github.com/skit-ai/vcore/surveillance"

	"github.com/getsentry/sentry-go"
	"github.com/google/go-cmp/cmp"
)

var DevMode = strings.ToUpper(env.String("ENVIRONMENT", "")) == "DEV"

var link = regexp.MustCompile("(^[A-Za-z])|_([A-Za-z])")

// StringInSlice - Returns True when two strings have one element in common, False otherwise
func StringInSlice(a string, list []string) bool {
	if list == nil {
		return false
	}
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// StringInSlice - Returns True when two strings have one element in common, False otherwise
func IntInSlice(a int, list []int) bool {
	if list == nil {
		return false
	}
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// SliceInSlice - Returns True when two []string have one element in common, False otherwise
func SliceInSlice(a []string, list []string) bool {
	for _, b := range list {
		if StringInSlice(b, a) {
			return true
		}
	}
	return false
}

// IsZeroOfUnderlyingType - Returns True when the passed value is zero
func IsZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

// IsDeepEqual - Returns True when the two passed values are same, else False
func IsDeepEqual(x, y interface{}) bool {
	return reflect.DeepEqual(x, y)
}

// IsEqual - Returns True when the two passed values are same, else False
func IsEqual(x, y interface{}) bool {
	return cmp.Equal(x, y)
}

// StringToTimestamp - Converts a unix timestamp string to unix timestamp format
func StringToTimestamp(unixTimestamp string) (time.Time, error) {
	epochFloat, err := strconv.ParseFloat(unixTimestamp, 64)
	if err != nil {
		log.Print(err)
		return time.Time{}, err
	}
	sec, dec := math.Modf(epochFloat)
	return time.Unix(int64(sec), int64(dec*(1e9))), nil
}

// Evaluate ...
func Evaluate(templateText string, metadata map[string]interface{}) string {
	return evaluate(templateText, metadata, nil)
}

// Recursively resolve all HTML template expressions
func evaluate(templateText string, metadata map[string]interface{}, funcMap template.FuncMap) string {
	text := templateText
	prevText := ""
	// Keep attempting to evaluate the expression as long as the no. of expressions reduces
	// This will result in the func recursively resolving all nested expressions until there are no more that can
	// be resolved
	var err error
	var tmpl *template.Template
	for text != prevText {
		prevText = text

		if funcMap != nil {
			tmpl, err = template.New("tts_dynamic").Funcs(funcMap).Parse(text)
		} else {
			tmpl, err = template.New("tts_dynamic").Parse(text)
		}

		if err != nil {
			errors.PrintStackTrace(err)
			break
		} else {
			var textBytes bytes.Buffer
			if err := tmpl.Execute(&textBytes, metadata); err != nil {
				errors.PrintStackTrace(err)
				break
			} else {
				text = textBytes.String()
			}
		}
	}

	return text
}

// Recursively resolve all HTML template expressions
func evaluateSafely(templateText string, metadata map[string]interface{}, funcMap template.FuncMap) (string, error) {
	text := templateText
	prevText := ""
	// Keep attempting to evaluate the expression as long as the no. of expressions reduces
	// This will result in the func recursively resolving all nested expressions until there are no more that can
	// be resolved
	var err error
	var tmpl *template.Template
	for text != prevText {
		prevText = text

		if funcMap != nil {
			tmpl, err = template.New("tts_dynamic").Funcs(funcMap).Parse(text)
		} else {
			tmpl, err = template.New("tts_dynamic").Parse(text)
		}

		if err != nil {
			return text, err
		} else {
			var textBytes bytes.Buffer
			if err := tmpl.Execute(&textBytes, metadata); err != nil {
				return text, err
			} else {
				text = textBytes.String()
			}
		}
	}

	return text, nil
}

func EvaluateAugmentedFuncMap(templateText string, metadata map[string]interface{}, additionalFuncs template.FuncMap) string {
	return evaluate(templateText, metadata, additionalFuncs)
}

func EvaluateAugmentedFuncMapSafely(templateText string, metadata map[string]interface{}, additionalFuncs template.FuncMap) (string, error) {
	return evaluateSafely(templateText, metadata, additionalFuncs)
}

// TrimSuffix - To trim a suffix string
func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

const (
	IgnoreRepeats = iota
	IgnoreOrdering
	ConsiderOrdering
)

// Checks if a string slice is equal or not.
func IsStringSliceEqual(a, b []string, option int) bool {
	switch option {
	case IgnoreRepeats:
		// Ignore repeated elements and treat both slices as a set and compare their elements
		aMap := buildKeyMap(a)
		bMap := buildKeyMap(b)
		for aMapKey, _ := range aMap {
			if _, exists := bMap[aMapKey]; !exists {
				return false
			}
		}
	case IgnoreOrdering:
		// If ordering is to be ignored, then we are to check if:
		// - All unique elements are the same
		// Checking equality of the elements regardless of the ordering
		// - All repeated elements are repeated exactly the same number of times in both slices
		xMap := make(map[string]int)
		yMap := make(map[string]int)
		for _, xElem := range a {
			xMap[xElem]++
		}
		for _, yElem := range a {
			yMap[yElem]++
		}

		for xMapKey, xMapVal := range xMap {
			if yMap[xMapKey] != xMapVal {
				return false
			}
		}
	case ConsiderOrdering:
		// simply loop through both loops and match elements
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
	}
	return true
}

func buildKeyMap(slice []string) map[string]struct{} {
	m := map[string]struct{}{}

	for _, v := range slice {
		m[v] = struct{}{}
	}
	return m
}

// Prints the time taken to run a function. Should be used to measure performance.
// Usage:
// Add  "defer TimeTrack(time.Now())" as a statement at the beginning of any function
func TimeTrack(start time.Time) {
	elapsed := time.Since(start)

	// Skip this function, and fetch the PC and file for its parent.
	pc, _, _, _ := runtime.Caller(1)

	// Retrieve a function object this functions parent.
	funcObj := runtime.FuncForPC(pc)

	log.Println(fmt.Sprintf("%s took %s", funcObj.Name(), elapsed))
}

// Gets a slice of distince strings from a slice of strings
func Distinct(slice []string) []string {
	var unique []string
	m := map[string]bool{}

	for _, v := range slice {
		if !m[v] {
			m[v] = true
			unique = append(unique, v)
		}
	}
	return unique
}

// Concatenates a slice of strings with the delimiter in question
// Uses the final delimiter to concatenate the final string in the slice
func StringListify(slice []string, delimiter, finalDelimiter string) string {
	var builder strings.Builder
	l := len(slice)
	for i, item := range slice {
		if i == l-1 && l > 1 {
			builder.WriteString(finalDelimiter)
		}
		builder.WriteString(item)
		if i < l-2 && l > 2 {
			builder.WriteString(delimiter)
		}
	}
	return builder.String()
}

// Concatenates a slice of strings with the delimiter in question
func Join(slice []string, delimiter string) string {
	var builder strings.Builder
	for i, item := range slice {
		builder.WriteString(item)
		if i != len(slice)-1 {
			builder.WriteString(delimiter)
		}
	}
	return builder.String()
}

// Concatenates a variable slice of strings
func Concat(slice ...interface{}) string {
	var builder strings.Builder
	for _, item := range slice {
		builder.WriteString(fmt.Sprintf("%v", item))
	}
	return builder.String()
}

// Concatenates a variable slice of strings
func JoinInt(delimiter string, slice []int) string {
	var builder strings.Builder
	for i, item := range slice {
		builder.WriteString(fmt.Sprintf("%v", item))
		if i != len(slice)-1 {
			builder.WriteString(delimiter)
		}
	}
	return builder.String()
}

// flattens a slice of []interface{} to a []string.
// The flattening is done only for one level of nesting currently
func FlattenStringSlice(nestedList []interface{}) (slice []string) {
	for _, element := range nestedList {
		switch item := element.(type) {
		case []string:
			slice = append(slice, item...)
		case string:
			slice = append(slice, item)
		case []interface{}:
			for _, item := range item {
				if val, ok := item.(string); ok {
					slice = append(slice, val)
				}
			}
		}
	}
	return
}

func ToCamelCase(str string) string {
	return link.ReplaceAllStringFunc(str, func(s string) string {
		return strings.ToUpper(strings.Replace(s, "_", " ", -1))
	})
}

// Closes a struct which implements io.Closer safely
func CloseSafely(closeable io.Closer) {
	if closeable != nil {
		surveillance.SentryClient.Capture(closeable.Close(), false)
	}
}

// Closes an interface which implements io.Closer safely
func CloseAnythingSafely(toClose interface{}) {
	if toClose != nil {
		if closeable, ok := toClose.(io.Closer); ok {
			CloseSafely(closeable)
		}
	}
}

// Handles an error by capturing it on Sentry and logging the same on STDOUT
func Capture(err error, _panic bool) sentry.EventID {
	return surveillance.SentryClient.Capture(err, _panic)
}

// Handles an error by capturing it on Sentry and logging the same on STDOUT
func CaptureWithContext(c context.Context, err error, _panic bool) sentry.EventID {
	return surveillance.SentryClient.CaptureWithContext(c, err, _panic)
}

func StringifyToJson(i interface{}) (stringifiedJson string) {
	if _formattedBytes, err := json.MarshalIndent(i, "", "\t"); err != nil {
		errors.PrintStackTrace(err)
	} else {
		stringifiedJson = string(_formattedBytes)
	}

	return
}

func ToString(toConvert interface{}) string {
	return fmt.Sprintf("%v", toConvert)
}

// Reverse - To reverse any string
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// MD5 - Generate MD5 hash for the given string
func MD5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

// ConvertInterfaceToString takes an interface as input and tries
// to assert it as a string if not nil. It returns the string value
// after the said assertion. False is returned if assertion fails
func ConvertInterfaceToString(val interface{}) (string, bool) {
	if val == nil {
		return "", false
	}

	stringValue, ok := val.(string)
	return stringValue, ok
}

// GetEnv - Return the fallback env if not present
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
