package ports

import "context"

type TxManager interface {
	WithinTx(context.Context, func(ctx context.Context) error) error
}