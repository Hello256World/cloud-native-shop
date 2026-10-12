package health

import "context"

type Checker interface {
	Name() string
	Check(context.Context) error
}

type checkerFunc struct {
	name string
	fn   func(context.Context) error
}

func NewChecker(name string, fn func(context.Context) error) Checker {
	return checkerFunc{name: name, fn: fn}
}

func (c checkerFunc) Name() string { return c.name }

func (c checkerFunc) Check(ctx context.Context) error { return c.fn(ctx) }
