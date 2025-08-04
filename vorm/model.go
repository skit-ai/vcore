package vorm

import (
	"database/sql/driver"
	"fmt"
	"github.com/skit-ai/vcore/errors"
	"os"
	"strconv"
	"time"

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

// Represents a database agnostic primary key
type Primary int64

func (x *Primary) Scan(src interface{}) error {
	if val, ok := src.(int64); ok {
		*x = Primary(val)
	} else if val, ok := src.(float64); ok {
		// This happens due to oracle driver not knowing that int64 is required by the calling code.
		*x = Primary(int64(val))
	} else if val, ok := src.(uint); ok {
		*x = Primary(int64(val))
	}
	// Note: This else results in errors while inserting into the table in oracle. Hence not using such a clause
	//else {
	//	return errors.New(fmt.Sprintf("Unable to convert %v to uint", src))
	//}

	return nil
}

func (x *Primary) Value() (driver.Value, error) {
	return *x, nil
}

func (x *Primary) Uint() uint {
	return uint(*x)
}

func (x *Primary) Foreign() Foreign {
	return Foreign(*x)
}

// Represents a database agnostic foreign key
type Foreign int64

func (x *Foreign) Scan(src interface{}) error {
	if val, ok := src.(uint); ok {
		*x = Foreign(int64(val))
	} else if val, ok := src.(float64); ok {
		// This happens due to oracle driver not knowing that uint is required by the calling code.
		*x = Foreign(int64(val))
	} else if val, ok := src.(int64); ok {
		*x = Foreign(val)
	}
	// Note: This else results in errors while inserting into the table in oracle. Hence not using such a clause
	//else {
	//	return errors.New(fmt.Sprintf("Unable to convert %v to uint", src))
	//}

	return nil
}

func (x *Foreign) Value() (driver.Value, error) {
	return *x, nil
}

func (x *Foreign) Uint() uint {
	return uint(*x)
}

func (x *Foreign) Int64() int64 {
	return int64(*x)
}

func (x *Foreign) Primary() Primary {
	return Primary(*x)
}

// Mimicks gorm.Model except ID is of type Primary
type ORM struct {
	ID        Primary `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
