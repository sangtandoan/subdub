package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/sangtandoan/subscription_tracker/internal/models"
)

type SessionRepo interface {
	GetSessionByID(ctx context.Context, id uuid.UUID) (*models.Session, error)
	CreateSession(ctx context.Context, arg *CreateSessionParams) (*models.Session, error)
	RevokeSession(ctx context.Context, id uuid.UUID) error
	DeleteSession(ctx context.Context, id uuid.UUID) error
}

type sessionRepo struct {
	db *sql.DB
}

func NewSessionRepo(db *sql.DB) *sessionRepo {
	return &sessionRepo{db}
}

func (repo *sessionRepo) GetSessionByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Session, error) {
	query := "SELECT id, user_email, refresh_token, is_revoked, created_at, expires_at FROM sessions WHERE id = $1"

	row := repo.db.QueryRowContext(ctx, query, id)

	var session models.Session
	err := row.Scan(
		&session.ID,
		&session.UserEmail,
		&session.RefreshToken,
		&session.IsRevoked,
		&session.CreatedAt,
		&session.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

type CreateSessionParams struct {
	ID           uuid.UUID
	RefreshToken string
	UserEmail    string
	ExpiresAt    time.Time
}

func (repo *sessionRepo) CreateSession(
	ctx context.Context,
	arg *CreateSessionParams,
) (*models.Session, error) {
	query := `
		INSERT INTO sessions (id, refresh_token, user_email, expires_at) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, user_email, refresh_token, is_revoked, created_at, expires_at
	`

	row := repo.db.QueryRowContext(
		ctx,
		query,
		arg.ID,
		arg.RefreshToken,
		arg.UserEmail,
		arg.ExpiresAt,
	)

	var session models.Session
	err := row.Scan(
		&session.ID,
		&session.UserEmail,
		&session.RefreshToken,
		&session.IsRevoked,
		&session.CreatedAt,
		&session.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (repo *sessionRepo) RevokeSession(ctx context.Context, id uuid.UUID) error {
	query := "UPDATE sessions SET is_revoked = true WHERE id = $1"

	_, err := repo.db.ExecContext(ctx, query, id)

	return err
}

func (repo *sessionRepo) DeleteSession(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM sessions WHERE id = $1"

	_, err := repo.db.ExecContext(ctx, query, id)

	return err
}
