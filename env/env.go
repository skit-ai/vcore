package env

import (
	"os"
	"strconv"
)

// Bool looks up for boolean env variables and returns it.
// This returns false if the key is not found or the value if non-boolean
func Bool(key string, fallback bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	parseBool, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parseBool
}

// String looks up for a string env variables and returns it.
func String(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return value
}

// Int looks up for a integer env variables and returns it.
func Int(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	parseInt, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parseInt
}

// Float looks up for a float64 env variables and returns it.
func Float(key string, fallback float64) float64 {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	parseFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fallback
	}

	return parseFloat
}
