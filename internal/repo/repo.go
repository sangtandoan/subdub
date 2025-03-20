package repo

import (
	"context"
	"database/sql"
	"time"
)

const QueryTimeOut = time.Second * 10

type Repo struct {
	User         UserRepo
	Subscription SubscriptionRepo
	Session      SessionRepo
	AuthProvider AuthProviderRepo
	TX           TX
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		User:         NewUserRepo(db),
		Subscription: NewSubsciptionRepo(db),
		Session:      NewSessionRepo(db),
		AuthProvider: NewAuthProviderRepo(db),
		TX:           NewTX(db),
	}
}

type Executor interface {
	QueryRowContext(ctx context.Context, query string, params ...any) *sql.Row
	ExecContext(ctx context.Context, query string, params ...any) (sql.Result, error)
}

type TX interface {
	WithTx(ctx context.Context, f func(txContext context.Context) error) error
}

type tx struct{ db *sql.DB }

func NewTX(db *sql.DB) *tx {
	return &tx{db}
}

type TxKey struct{}

func (t *tx) WithTx(ctx context.Context, f func(txContext context.Context) error) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	txContext := context.WithValue(ctx, TxKey{}, tx)
	err = f(txContext)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			panic(err)
		}
	}

	return tx.Commit()
}

func getExcutor(ctx context.Context, db *sql.DB) Executor {
	tx := ctx.Value(TxKey{})
	if tx == nil {
		return db
	}

	return tx.(*sql.Tx)
}
