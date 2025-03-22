package repo

import (
	"context"
	"database/sql"
	"time"
)

// QueryTimeOut defines the standard timeout for database operations
const QueryTimeOut = time.Second * 10

// Repo contains all repository interfaces
type Repo struct {
	User         UserRepo
	Subscription SubscriptionRepo
	Session      SessionRepo
	AuthProvider AuthProviderRepo
	Transaction  TransactionManager
}

// NewRepo creates a new repository instance with all dependencies
func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		User:         NewUserRepo(db),
		Subscription: NewSubsciptionRepo(db),
		Session:      NewSessionRepo(db),
		AuthProvider: NewAuthProviderRepo(db),
		Transaction:  NewTransactionManager(db),
	}
}

type Executor interface {
	QueryRowContext(ctx context.Context, query string, params ...any) *sql.Row
	ExecContext(ctx context.Context, query string, params ...any) (sql.Result, error)
}

// TransactionManager defines the interface for transaction operations
type TransactionManager interface {
	WithTx(ctx context.Context, f func(txContext context.Context) error) error
}

type transactionManager struct{ db *sql.DB }

func NewTransactionManager(db *sql.DB) *transactionManager {
	return &transactionManager{db}
}

// TxKey is the context key for transaction
type TxKey struct{}

// WithTx executes the given function within a transaction
func (t *transactionManager) WithTx(
	ctx context.Context,
	f func(txContext context.Context) error,
) error {
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

// getExcutor returns the appropriate executor (DB or Transaction) from context
func getExcutor(ctx context.Context, db *sql.DB) Executor {
	tx, ok := ctx.Value(TxKey{}).(*sql.Tx)
	if !ok || tx == nil {
		return db
	}

	return tx
}
