package vorm

import (
	"fmt"
	"github.com/Vernacular-ai/vcore/errors"
	"os"
	"strconv"

	"github.com/Vernacular-ai/gorm"
	_ "github.com/Vernacular-ai/gorm/dialects/mysql"
	_ "github.com/Vernacular-ai/gorm/dialects/oci8"
	_ "github.com/Vernacular-ai/gorm/dialects/postgres"
)

const (
	PostgresDriver = "postgres"
	OracleDriver   = "oci8"
	MySQLDriver    = "mysql"
)

type Model struct {
	*gorm.DB
}

var DB *Model

// Connects to the DB based on the driver set in the env variable "DB_DRIVER"
// Currently supported values for DB_DRIVER = "postgres" for postgres, "oci8" for oracle
func InitDB(dataSourceName string) (*Model, error) {
	driver := os.Getenv("DB_DRIVER")
	switch driver {
	case PostgresDriver:
		return InitPostgresDB(dataSourceName)
	case OracleDriver:
		return InitOracleDB(dataSourceName)
	case MySQLDriver:
		return InitMySQLDB(dataSourceName)
	default:
		return nil, errors.NewError(fmt.Sprintf("Driver `%s` not supported.", driver), nil, true)
	}
}

func InitPostgresDB(dataSourceName string) (*Model, error) {
	return initDBInternal(PostgresDriver, dataSourceName)
}

func InitOracleDB(dataSourceName string) (*Model, error) {
	return initDBInternal(OracleDriver, dataSourceName)
}

func InitMySQLDB(dataSourceName string) (*Model, error) {
	return initDBInternal(MySQLDriver, dataSourceName)
}

func initDBInternal(dialect string, dataSourceName string) (*Model, error) {
	if db, err := gorm.Open(dialect, dataSourceName); err == nil {
		trace, _ := strconv.ParseBool(os.Getenv("DB_TRACE"))
		db.LogMode(trace)
		db.DB().SetMaxIdleConns(10)
		db.DB().SetMaxOpenConns(100)
		// Setting the DB in VORM as well
		instance := &Model{db}
		DB = instance
		return instance, nil
	} else {
		return nil, errors.NewError(fmt.Sprintf("Unable to connect using dialect `%s` and source `%s`", dialect, dataSourceName), err, true)
	}
}
