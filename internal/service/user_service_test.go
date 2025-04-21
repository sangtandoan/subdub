package service_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/sangtandoan/subscription_tracker/internal/service"
	"github.com/sangtandoan/subscription_tracker/internal/utils"
	"github.com/sangtandoan/subscription_tracker/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetUser(t *testing.T) {
	user := utils.RandomUser()

	testCases := []struct {
		buildStubs    func(*mocks.MockUserRepo)
		checkResponse func(*testing.T, *service.GetUserResponse, error)
		name          string
		userID        uuid.UUID
	}{
		{
			name:   "Get user successfully",
			userID: user.ID,
			buildStubs: func(repo *mocks.MockUserRepo) {
				repo.EXPECT().GetUserByID(gomock.Any(), user.ID).Times(1).Return(user, nil)
			},
			checkResponse: func(t *testing.T, response *service.GetUserResponse, err error) {
				require.NoError(t, err)

				require.NotEmpty(t, response)

				require.Equal(t, user.ID, response.ID)
				require.Equal(t, user.Email, response.Email)
				require.Equal(t, user.CreatedAt, response.CreatedAt)
			},
		},
		{
			name:   "No user found",
			userID: uuid.New(),
			buildStubs: func(repo *mocks.MockUserRepo) {
				repo.EXPECT().
					GetUserByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, response *service.GetUserResponse, err error) {
				require.Empty(t, response)
				require.Error(t, err)
				require.ErrorIs(t, err, sql.ErrNoRows)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockUserRepo(ctrl)
			tc.buildStubs(mockRepo)
			userService := service.NewUserService(mockRepo)

			response, err := userService.GetUser(context.Background(), tc.userID)

			tc.checkResponse(t, response, err)
		})
	}
}
