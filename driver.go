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
	if err != nil {
		return nil, err
	}
	if _, is := conn.(driver.ConnBeginTx); !is {
		panic("你的mysql driver没有实现driver.ConnBeginTx")
	}
	if _, is := conn.(driver.ConnPrepareContext); !is {
		panic("你的mysql driver没有实现driver.ConnPrepareContext")
	}
	if _, is := conn.(driver.ExecerContext); !is {
		panic("你的mysql driver没有实现driver.ExecerContext")
	}
	if _, is := conn.(driver.QueryerContext); !is {
		panic("你的mysql driver没有实现driver.QueryerContext")
	}
	return &zpConn{conn: conn.(baseConn)}, nil
}

var (
	_ driver.Driver = ZpMySQLDriver{}
)
