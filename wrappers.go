package zpmysql

import (
	"context"
	"database/sql/driver"
)

type connWithFakePrepareContext struct{ driver.Conn }

func (c *connWithFakePrepareContext) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	return c.Prepare(query)
}

type connWithFakeBeginTx struct{ driver.Conn }

func (c *connWithFakeBeginTx) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return c.Begin()
}

var (
	_ driver.ConnPrepareContext = &connWithFakePrepareContext{}
	_ driver.ConnBeginTx        = &connWithFakeBeginTx{}
)
