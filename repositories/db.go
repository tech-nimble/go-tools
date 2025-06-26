package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tech-nimble/go-tools/helpers/errors"
)

const TransactionKey = "pgx_transaction"

type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Repositories struct {
	DB *pgxpool.Pool
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	rep := &Repositories{
		DB: db,
	}

	return rep
}

func (r *Repositories) StartTransaction(ctx context.Context) (context.Context, error) {
	_, ok := r.GetTx(ctx)
	if ok {
		return ctx, nil
	}

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, TransactionKey, tx), nil //nolint
}

func (r *Repositories) StopTransaction(ctx context.Context) error {
	tx, ok := r.GetTx(ctx)
	if !ok {
		return errors.Runtime.New("Transaction not started ")
	}

	return tx.Commit(ctx)
}

func (r *Repositories) CancelTransaction(ctx context.Context) error {
	tx, ok := r.GetTx(ctx)
	if !ok {
		return errors.Runtime.New("Transaction not started ")
	}

	return tx.Rollback(ctx)
}

func (r *Repositories) GetTx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(TransactionKey).(pgx.Tx)

	return tx, ok
}

func (r *Repositories) GetConnect(ctx context.Context) Querier {
	tx, ok := r.GetTx(ctx)
	if !ok {
		return r.DB
	}

	return tx
}
