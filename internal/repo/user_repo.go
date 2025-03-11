package repo

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/sangtandoan/subscription_tracker/internal/models"
)

type UserRepo interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, arg *CreateUserParams) (*models.User, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *userRepo {
	return &userRepo{db}
}

func (repo *userRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := "SELECT id, email, password, created_at FROM users WHERE id = $1"

	timeOutCtx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	row := repo.db.QueryRowContext(timeOutCtx, query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *userRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := "SELECT id, email, password, created_at FROM users WHERE email = $1"

	timeOutCtx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	row := repo.db.QueryRowContext(timeOutCtx, query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

type CreateUserParams struct {
	ID       uuid.UUID
	Email    string
	Password string
}

func (repo *userRepo) CreateUser(ctx context.Context, arg *CreateUserParams) (*models.User, error) {
	query := "INSERT INTO users (id, email, password) VALUES ($1, $2, $3) RETURNING id, email, password, created_at"

	timeOutCtx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	row := repo.db.QueryRowContext(timeOutCtx, query, arg.ID, arg.Email, arg.Password)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
