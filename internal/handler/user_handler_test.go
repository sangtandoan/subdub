package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sangtandoan/subscription_tracker/internal/authenticator"
	"github.com/sangtandoan/subscription_tracker/internal/handler"
	"github.com/sangtandoan/subscription_tracker/internal/models"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/response"
	"github.com/sangtandoan/subscription_tracker/internal/service"
	"github.com/sangtandoan/subscription_tracker/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetUserHandler(t *testing.T) {
	user := randomUser()
	getUserResponse := &service.GetUserResponse{
		CreatedAt: user.CreatedAt,
		Email:     user.Email,
		ID:        user.ID,
	}
	appResponseJSON, _ := json.Marshal(
		response.NewAppResponse("get user successfully", getUserResponse),
	)

	testCases := []struct {
		buildStubs    func(*mocks.MockUserService)
		checkResponse func(*testing.T, *httptest.ResponseRecorder)
		name          string
		userID        string
	}{
		{
			name:   "Get user successfully",
			userID: user.ID.String(),
			buildStubs: func(s *mocks.MockUserService) {
				s.EXPECT().GetUser(gomock.Any(), user.ID).Times(1).Return(getUserResponse, nil)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, rec.Code, http.StatusOK)
				require.JSONEq(
					t,
					string(appResponseJSON),
					rec.Body.String(),
				)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gin.SetMode(gin.TestMode)

			userService := mocks.NewMockUserService(ctrl)
			tc.buildStubs(userService)
			userHandler := handler.NewUserHandler(userService)

			// act
			req := httptest.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/api/v1/users/%s", tc.userID),
				nil,
			)
			rec := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(rec)
			c.Set(authenticator.SubClaim, tc.userID)
			c.Params = gin.Params{
				{Key: "id", Value: tc.userID},
			}
			c.Request = req

			userHandler.GetUserHandler(c)
			// assert
			tc.checkResponse(t, rec)
		})
	}
}

func randomUser() *models.User {
	return &models.User{
		CreatedAt: time.Now(),
		Email:     "sangvaminh11497@gmail.com",
		Password:  "secret",
		ID:        uuid.New(),
	}
}
