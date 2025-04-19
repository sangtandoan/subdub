package repo

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/sangtandoan/subscription_tracker/internal/models"
)

type AuthProviderRepo interface {
	CreateAuthProvider(
		ctx context.Context,
		arg *CreateAuthProviderParams,
	) (*models.AuthProvider, error)
}

type authProviderRepo struct {
	db *sql.DB
}

func NewAuthProviderRepo(db *sql.DB) *authProviderRepo {
	return &authProviderRepo{db}
}

type CreateAuthProviderParams struct {
	Provider   string
	ProviderID string
	ID         uuid.UUID
	UserID     uuid.UUID
}

func (repo *authProviderRepo) CreateAuthProvider(
	ctx context.Context,
	arg *CreateAuthProviderParams,
) (*models.AuthProvider, error) {
	ex := getExcutor(ctx, repo.db)

	query := `
		INSERT INTO auth_providers (id, user_id, provider, provider_id) 
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, provider, provider_id, created_at`

	row := ex.QueryRowContext(ctx, query, arg.ID, arg.UserID, arg.Provider, arg.ProviderID)

	var authProvider models.AuthProvider
	err := row.Scan(
		&authProvider.ID,
		&authProvider.UserID,
		&authProvider.Provider,
		&authProvider.ProviderID,
		&authProvider.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &authProvider, nil
}
