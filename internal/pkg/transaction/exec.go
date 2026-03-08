package transaction

import (
	"context"

	"fg_bot/internal/pkg/transaction/abstract"
)

func Exec(ctx context.Context, callback func(ctx context.Context) error) (err error) {
	var commit abstract.CommitFunc
	ctx, commit = abstract.ContextWithTx(ctx, nil)
	defer func() {
		err = commit(err)
	}()
	return callback(ctx)
}
