package vorm

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"github.com/Vernacular-ai/vcore/errors"
	"github.com/Vernacular-ai/vcore/surveillance"
)

// Returns the formatted value of the JSON field based on the expectations of the target driver
func jsonValue(i interface{}) (driver.Value, error) {
	var err error
	if DB == nil {
		return nil, errors.NewError("DB connection has not been initialized", nil, true)
	} else {
		// Determing the dialect
		dialect := DB.Dialect().GetName()

		// Convert the jsonValue into bytes
		var b []byte
		if b, err = json.Marshal(i); err != nil {
			return nil, errors.NewError("Unable to marshal JSON", err, false)
		}

		switch dialect {
		case OracleDriver:
			// Requires string instead of bytes
			return string(b), nil
		case PostgresDriver:
			// Requires bytes
			return b, nil
		default:
			return nil, errors.NewError(fmt.Sprintf("Dialect `%s` not supported.", dialect), nil, true)
		}
	}
}

// Use in migrations in case the DB entry is to be made in JSON format
// Can use JsonbMap or StringListMap to automatically unmarshal from db into the desired data type in the actual code
type ORMJson struct {
	json.RawMessage
}

func (j ORMJson) Value() (val driver.Value, err error) {
	return jsonValue(j)
}

func (j *ORMJson) Scan(src interface{}) error {
	switch source := src.(type) {
	case string:
		// Receives interface as a string
		*j = ORMJson{json.RawMessage([]byte(source))}
	case []byte:
		// Receives a [] bytes
		*j = ORMJson{json.RawMessage(source)}
	default:
		return errors.NewError(fmt.Sprintf("Type `%s` not supported.", reflect.TypeOf(source)), nil, true)
	}
	return nil
}

// Creates a postgres JSON from an interface
func CreateORMJson(i interface{}) ORMJson {
	b, err := json.Marshal(i)
	if err != nil {
		surveillance.SentryClient.Capture(err, false)
		return ORMJson{}
	} else {
		return ORMJson{json.RawMessage(b)}
	}
}

// Represents a string list stored in the DB in JSON format.
type JsonStringSlice []string

func (l *JsonStringSlice) Value() (val driver.Value, err error) {
	return jsonValue(l)
}

func (l *JsonStringSlice) Scan(src interface{}) error {
	switch source := src.(type) {
	case string:
		// Receives interface as a string
		return l.setValue([]byte(source))
	case []byte:
		// Receives a [] bytes
		return l.setValue(source)
	default:
		return errors.NewError(fmt.Sprintf("Type `%s` not supported.", reflect.TypeOf(source)), nil, true)
	}
}

// Convert bytes to a []string
func (l *JsonStringSlice) setValue(bytes []byte) (err error) {
	var i []string
	if err = json.Unmarshal(bytes, &i); err != nil {
		err = errors.NewError("Unable to unmarshal JSON into ([]string) ", err, false)
	} else {
		*l = i
	}
	return
}

// Automatically converts JSON type from DB into a map[string]interface{}
type JsonbMap map[string]interface{}

func (p JsonbMap) Value() (val driver.Value, err error) {
	return jsonValue(p)
}

func (p *JsonbMap) Scan(src interface{}) error {
	switch source := src.(type) {
	case string:
		// Receives interface as a string
		return p.setValue([]byte(source))
	case []byte:
		// Receives a [] bytes
		return p.setValue(source)
	default:
		return errors.NewError(fmt.Sprintf("Type `%s` not supported.", reflect.TypeOf(source)), nil, true)
	}
}

// Convert bytes to a map[string]interface{}
func (p *JsonbMap) setValue(bytes []byte) (err error) {
	var i map[string]interface{}
	if err = json.Unmarshal(bytes, &i); err != nil {
		err = errors.NewError("Unable to unmarshal JSON into (map[string]interface{}) ", err, false)
	} else {
		*p = i
	}
	return
}
