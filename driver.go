package zpmysql

import (
	"database/sql/driver"
)

type ZpMySQLDriver struct {
	driver driver.Driver
}

func (d ZpMySQLDriver) Open(dsn string) (driver.Conn, error) {
	conn, err := d.driver.Open(dsn)
	AfterProcess(&HookContext{Err: err})
	return &zpConn{conn: conn}, err
}

var (
	_ driver.Driver = ZpMySQLDriver{}
)
