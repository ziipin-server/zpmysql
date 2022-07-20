package zpmysql

import "database/sql/driver"

type zpTx struct {
	tx driver.Tx
}

func (tx *zpTx) Commit() error {
	return AfterProcess(&HookContext{Err: tx.tx.Commit()})
}

func (tx *zpTx) Rollback() error {
	return AfterProcess(&HookContext{Err: tx.tx.Rollback()})
}

var (
	_ driver.Tx = &zpTx{}
)
