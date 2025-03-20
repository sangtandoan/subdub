package repo

import (
	"context"
	"database/sql"
	"time"

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

	return scanUser(row)
}

func (repo *userRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := "SELECT id, email, password, created_at FROM users WHERE email = $1"

	timeOutCtx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	row := repo.db.QueryRowContext(timeOutCtx, query, email)

	return scanUser(row)
}

type CreateUserParams struct {
	ID       uuid.UUID
	Email    string
	Password string
}

type UserRow struct {
	Password  *string
	CreatedAt time.Time
	ID        uuid.UUID
	Email     string
}

func (repo *userRepo) CreateUser(ctx context.Context, arg *CreateUserParams) (*models.User, error) {
	ex := getExcutor(ctx, repo.db)

	query := "INSERT INTO users (id, email, password) VALUES ($1, $2, $3) RETURNING id, email, password, created_at"

	if arg.Password == "" {
		query = "INSERT INTO users (id, email) VALUES ($1, $2) RETURNING id, email, password, created_at"
	}

	timeOutCtx, cancel := context.WithTimeout(ctx, QueryTimeOut)
	defer cancel()

	var row *sql.Row
	if arg.Password != "" {
		row = ex.QueryRowContext(timeOutCtx, query, arg.ID, arg.Email, arg.Password)
	} else {
		row = ex.QueryRowContext(timeOutCtx, query, arg.ID, arg.Email)
	}

	return scanUser(row)
}

func toUser(user *UserRow, password string) *models.User {
	return &models.User{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		Password:  password,
	}
}

func scanUser(row *sql.Row) (*models.User, error) {
	var user UserRow
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	password := ""
	if user.Password != nil {
		password = *user.Password
	}

	return toUser(&user, password), nil
}
