package connector

import (
	"fmt"

	"dbmw/connector/mariadb"
	"dbmw/connector/mysql"
	"dbmw/connector/postgres"
	"dbmw/connector/sqlite"
	"dbmw/core/connection"
)

// DefaultFactory creates a new uninitialized Connector for any supported driver.
func DefaultFactory(driver connection.DriverType) (connection.Connector, error) {
	switch driver {
	case connection.DriverPostgres:
		return postgres.NewConnector(), nil
	case connection.DriverMySQL:
		return mysql.NewConnector(), nil
	case connection.DriverMariaDB:
		return mariadb.NewConnector(), nil
	case connection.DriverSQLite:
		return sqlite.NewConnector(), nil
	default:
		return nil, fmt.Errorf("%w: %s", connection.ErrUnsupportedDriver, driver)
	}
}
