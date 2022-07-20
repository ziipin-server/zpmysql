package zpmysql

import (
	"context"
	"database/sql/driver"
)

type zpStmt struct {
	stmt driver.Stmt
	sql  string
}

func (s *zpStmt) Close() error {
	err := s.stmt.Close()
	if err != nil {
		AfterProcess(&HookContext{Err: err})
	}
	return err
}

func (s *zpStmt) NumInput() int {
	return s.stmt.NumInput()
}

// Deprecated: 兼容
func (s *zpStmt) Exec(args []driver.Value) (driver.Result, error) {
	result, err := s.stmt.Exec(args)
	AfterProcess(&HookContext{Err: err})
	return result, err
}

// Deprecated: 兼容
func (s *zpStmt) Query(args []driver.Value) (driver.Rows, error) {
	result, err := s.stmt.Query(args)
	AfterProcess(&HookContext{Err: err})
	return result, err
}

type zpStmtExecContext struct {
	zpStmt
}

func (s *zpStmtExecContext) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	hookCtx := NewHookContext(ctx, s.zpStmt.sql, args)
	ctx, _ = BeforeProcess(hookCtx)
	r, err := s.zpStmt.stmt.(driver.StmtExecContext).ExecContext(ctx, args)
	hookCtx.End(ctx, err)
	AfterProcess(hookCtx)
	return r, err
}

type zpStmtQueryContext struct {
	zpStmt
}

func (s *zpStmtQueryContext) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	hookCtx := NewHookContext(ctx, s.zpStmt.sql, args)
	ctx, _ = BeforeProcess(hookCtx)
	r, err := s.zpStmt.stmt.(driver.StmtQueryContext).QueryContext(ctx, args)
	hookCtx.End(ctx, err)
	AfterProcess(hookCtx)
	return r, err
}

type zpStmtExecQueryContext struct {
	zpStmt
}

func (s *zpStmtExecQueryContext) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	hookCtx := NewHookContext(ctx, s.zpStmt.sql, args)
	ctx, _ = BeforeProcess(hookCtx)
	r, err := s.zpStmt.stmt.(driver.StmtExecContext).ExecContext(ctx, args)
	hookCtx.End(ctx, err)
	AfterProcess(hookCtx)
	return r, err
}

func (s *zpStmtExecQueryContext) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	hookCtx := NewHookContext(ctx, s.zpStmt.sql, args)
	ctx, _ = BeforeProcess(hookCtx)
	r, err := s.zpStmt.stmt.(driver.StmtQueryContext).QueryContext(ctx, args)
	hookCtx.End(ctx, err)
	AfterProcess(hookCtx)
	return r, err
}

func newStmt(stmt driver.Stmt, sql string) driver.Stmt {
	switch stmt.(type) {
	case driver.StmtExecContext:
		if _, ok := stmt.(driver.StmtQueryContext); ok {
			return &zpStmtExecQueryContext{zpStmt: zpStmt{stmt: stmt, sql: sql}}
		} else {
			return &zpStmtExecContext{zpStmt: zpStmt{stmt: stmt, sql: sql}}
		}
	case driver.StmtQueryContext:
		return &zpStmtQueryContext{zpStmt: zpStmt{stmt: stmt, sql: sql}}
	default:
		return &zpStmt{stmt: stmt, sql: sql}
	}
}

var (
	_ driver.Stmt            = &zpStmt{}
	_ driver.StmtExecContext = &zpStmtExecContext{}
	_ driver.Stmt            = &zpStmtQueryContext{}
	_ driver.Stmt            = &zpStmtExecQueryContext{}
)
