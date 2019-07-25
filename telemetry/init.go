package telemetry

import (
	"net/http"
	"os"

	"github.com/newrelic/go-agent"
)

// Wrapper for newrelic.Application
type NewRelicClient struct {
	application newrelic.Application
}

// Wrapper for newrelic.Transaction
type Transaction struct {
	transaction newrelic.Transaction
}

func newClient() NewRelicClient{
	appName := os.Getenv("NEW_RELIC_APP_NAME")
	licenseKey := os.Getenv("NEW_RELIC_LICENSE_KEY")

	if appName != "" && licenseKey != "" {
		// Initialize Newrelic
		config := newrelic.NewConfig(appName, licenseKey)
		if application, err := newrelic.NewApplication(config); err == nil {
			// Only initialize the newrelic client if there is no error
			return NewRelicClient{application}
		} else {
			return NewRelicClient{nil}
		}
	} else {
		return NewRelicClient{nil}
	}
}

var (
	//Global Newrelic Client
	NewrelicClient = newClient()
)

// Wrapper to only start the transaction if the newrelic application struct was successfully instantiated
func (c *NewRelicClient) StartTransaction(name string, w http.ResponseWriter, r *http.Request) Transaction {
	if c.application != nil {
		return Transaction{c.application.StartTransaction(name, w, r)}
	} else {
		return Transaction{nil}
	}
}

// Wrapper to end a transaction only if the newrelic transaction struct was created
func (t *Transaction) End() error {
	if t.transaction != nil {
		return t.transaction.End()
	} else {
		return nil
	}
}
