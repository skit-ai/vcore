package telemetry

import (
	"os"

	newrelic "github.com/newrelic/go-agent"
)

var (
	//NewrelicClient
	NewrelicClient newrelic.Application
)

func Init() {

	// Initialize Newrelic
	config := newrelic.NewConfig(os.Getenv("NEW_RELIC_APP_NAME"), os.Getenv("NEW_RELIC_LICENSE_KEY"))
	NewrelicClient, _ = newrelic.NewApplication(config)
}
