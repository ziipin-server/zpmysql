package zpmysql

import (
	"context"
	"database/sql/driver"
)

type zpConn struct {
	conn driver.Conn
}

func (c *zpConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	hookCtx := NewHookContext(ctx, "PREPARE", nil)
	ctx, _ = BeforeProcess(hookCtx)
	stmt, err := c.conn.Prepare(query)
	hookCtx.End(ctx, err)
	AfterProcess(hookCtx)
	return newStmt(stmt, query), err
}

func (c *zpConn) Prepare(query string) (driver.Stmt, error) {
	stmt, err := c.conn.Prepare(query)
	AfterProcess(&HookContext{Err: err})
	return newStmt(stmt, query), err
}

func (c *zpConn) Close() error {
	return AfterProcess(&HookContext{Err: c.conn.Close()})
}

// Deprecated: 兼容
func (c *zpConn) Begin() (driver.Tx, error) {
	tx, err := c.conn.Begin()

	return &zpTx{tx: tx}, err
}

var (
	_ driver.Conn = &zpConn{}
)
