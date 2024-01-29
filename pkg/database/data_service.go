package database

import "context"

type TxKey struct{}

func injectTx(ctx context.Context, s SqlGdbc) context.Context {
	return context.WithValue(ctx, TxKey{}, s)
}

func extractTx(ctx context.Context) SqlGdbc {
	if tx, ok := ctx.Value(TxKey{}).(SqlGdbc); ok {
		return tx
	}
	return nil
}
