package zpmysql

import (
	"context"
	"database/sql/driver"
)

type baseConn interface {
	driver.Conn
	driver.ConnPrepareContext
	driver.ConnBeginTx
	driver.ExecerContext
	driver.QueryerContext
}

type zpConn struct {
	conn baseConn
}

func (c *zpConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	hookCtx := NewHookContext(ctx, "PREPARE", nil)
	ctx, _ = BeforeProcess(hookCtx)
	stmt, err := c.conn.PrepareContext(ctx, query)
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
	AfterProcess(&HookContext{Err: err})
	return &zpTx{tx: tx}, err
}

func (c *zpConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	hookCtx := NewHookContext(ctx, "BEGIN", nil)
	tx, err := c.conn.BeginTx(ctx, opts)
	hookCtx.End(ctx, err)
	AfterProcess(hookCtx)
	return &zpTx{tx: tx}, err
}

func (c *zpConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	hookCtx := NewHookContext(ctx, query, nil)
	result, err := c.conn.ExecContext(ctx, query, args)
	hookCtx.End(ctx, err)
	AfterProcess(hookCtx)
	return result, err
}

func (c *zpConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	hookCtx := NewHookContext(ctx, query, nil)
	result, err := c.conn.QueryContext(ctx, query, args)
	hookCtx.End(ctx, err)
	AfterProcess(hookCtx)
	return result, err
}

var (
	_ baseConn = &zpConn{}
)
