package zpmysql

import (
	"context"
	"database/sql/driver"
)

type HookContext struct {
	Ctx  context.Context
	SQL  string
	Args []driver.NamedValue
	Err  error
}

func NewHookContext(ctx context.Context, sql string, args []driver.NamedValue) *HookContext {
	return &HookContext{
		Ctx:  ctx,
		SQL:  sql,
		Args: args,
	}
}

func (c *HookContext) End(ctx context.Context, err error) {
	c.Ctx = ctx
	c.Err = err
}

type Hook interface {
	BeforeProcess(c *HookContext) (context.Context, error)
	AfterProcess(c *HookContext) error
}

func BeforeProcess(c *HookContext) (context.Context, error) {
	ctx := c.Ctx
	for _, h := range hooks {
		var err error
		ctx, err = h.BeforeProcess(c)
		if err != nil {
			return nil, err
		}
	}
	return ctx, nil
}

func AfterProcess(c *HookContext) error {
	firstErr := c.Err
	for _, h := range hooks {
		err := h.AfterProcess(c)
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func AddHook(hook Hook) {
	hooks = append(hooks, hook)
}

var (
	hooks []Hook
)

func init() {
	hooks = make([]Hook, 0, 4)
}
