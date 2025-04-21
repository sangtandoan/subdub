package repo_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/sangtandoan/subscription_tracker/internal/models"
	"github.com/sangtandoan/subscription_tracker/internal/repo"
	"github.com/sangtandoan/subscription_tracker/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestGetUserByID(t *testing.T) {
	// db is mock for *sql.DB, mock is using for expectations for that *sql.DB
	//
	// It is different from using gomock because gomock creates a mock for that interface and
	// also uses that mock for expectations.
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	userID := uuid.New()
	user := utils.RandomUser()
	rows := randomRow(user)

	testCases := []struct {
		buildStubs    func(sqlmock.Sqlmock)
		checkResponse func(*testing.T, *models.User, error)
		name          string
		userID        uuid.UUID
	}{
		{
			name:   "Get user by ID successfully",
			userID: userID,
			buildStubs: func(mock sqlmock.Sqlmock) {
				query := `SELECT id, email, password, created_at FROM users WHERE id = \$1`
				// ExpectQuery need regex string to match the query needed to be tested
				mock.ExpectQuery(query).WithArgs(userID).WillReturnRows(rows)
			},
			checkResponse: func(t *testing.T, response *models.User, err error) {
				require.NoError(t, err)

				require.NotNil(t, response)

				require.Equal(t, user.ID, response.ID)
				require.Equal(t, user.Email, response.Email)
				require.Equal(t, user.Password, response.Password)
				require.Equal(t, user.CreatedAt, response.CreatedAt)
			},
		},
		{
			name:   "Get user by ID not found",
			userID: userID,
			buildStubs: func(mock sqlmock.Sqlmock) {
				query := `SELECT id, email, password, created_at FROM users WHERE id = \$1`
				mock.ExpectQuery(query).WithArgs(userID).WillReturnError(sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, response *models.User, err error) {
				require.Error(t, err)

				require.ErrorIs(t, err, sql.ErrNoRows)

				require.Nil(t, response)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(mock)

			repo := repo.NewUserRepo(db)
			response, err := repo.GetUserByID(context.Background(), tc.userID)

			tc.checkResponse(t, response, err)

			// Make sure all expectations were met
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func randomRow(user *models.User) *sqlmock.Rows {
	// AddRows will return *sql.Rows for array of rows
	// AddRow will return *sql.Row for single row
	return sqlmock.NewRows([]string{"id", "email", "password", "created_at"}).
		AddRow(user.ID, user.Email, user.Password, user.CreatedAt)
}
