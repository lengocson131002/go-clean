package database

import (
	"context"
	"fmt"
)

func (sdt *SqlxDBTx) WithinTransaction(ctx context.Context, txFunc func(ctx context.Context) error) error {
	var err error
	tx, err := sdt.DB.Beginx()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	sct := &SqlxConnTx{tx}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = txFunc(InjectTx(ctx, sct))
	return err
}

func (sct *SqlxConnTx) WithinTransaction(ctx context.Context, txFunc func(ctx context.Context) error) error {
	var err error
	tx := sct.DB
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = txFunc(InjectTx(ctx, sct))
	return err
}
