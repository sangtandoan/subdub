package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sangtandoan/subscription_tracker/internal/authenticator"
	"github.com/sangtandoan/subscription_tracker/internal/handler"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/apperror"
	"github.com/sangtandoan/subscription_tracker/internal/pkg/response"
	"github.com/sangtandoan/subscription_tracker/internal/service"
	"github.com/sangtandoan/subscription_tracker/internal/utils"
	"github.com/sangtandoan/subscription_tracker/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetUserHandler(t *testing.T) {
	user := utils.RandomUser()
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
		checkResponse func(*testing.T, *gin.Context, *httptest.ResponseRecorder)
		setContext    func(*gin.Context, string)
		name          string
		userID        string
	}{
		{
			name:   "Get user successfully",
			userID: user.ID.String(),
			setContext: func(c *gin.Context, userID string) {
				c.Set(authenticator.SubClaim, userID)

				// We need to set c.Params here even the URL has the id because
				// we call handler directly without going through the gin router,
				// where the URL parameters are automatically parsed and set in c.Params.
				c.Params = gin.Params{
					{Key: "id", Value: userID},
				}
			},
			buildStubs: func(s *mocks.MockUserService) {
				s.EXPECT().GetUser(gomock.Any(), user.ID).Times(1).Return(getUserResponse, nil)
			},
			checkResponse: func(t *testing.T, c *gin.Context, rec *httptest.ResponseRecorder) {
				require.Equal(t, rec.Code, http.StatusOK)
				require.JSONEq(
					t,
					string(appResponseJSON),
					rec.Body.String(),
				)
			},
		},
		{
			name:   "Get user with invalid id",
			userID: "invalid-uuid",
			setContext: func(c *gin.Context, userID string) {
			},
			buildStubs: func(s *mocks.MockUserService) {
				s.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, c *gin.Context, rec *httptest.ResponseRecorder) {
				require.Greater(t, len(c.Errors), 0)

				for _, err := range c.Errors {
					// ErrorIs is used to check if at least one of the errors in the error's chain
					// is of a specific type
					// ErrorAs is used to check if at least one of the errors in the error's chain
					// is of a specific type and then ssigned to a variable of that type
					require.ErrorIs(t, err.Err, apperror.ErrInvalidJSON)
				}
			},
		},
		{
			name:   "Get user without userID in context",
			userID: user.ID.String(),
			setContext: func(c *gin.Context, userID string) {
				c.Params = gin.Params{
					{Key: "id", Value: userID},
				}
			},
			buildStubs: func(s *mocks.MockUserService) {
				s.EXPECT().GetUser(gomock.Any(), user.ID).Times(0)
			},
			checkResponse: func(t *testing.T, c *gin.Context, rec *httptest.ResponseRecorder) {
				require.Greater(t, len(c.Errors), 0)

				for _, err := range c.Errors {
					require.ErrorIs(t, err.Err, apperror.ErrUnAuthorized)
				}
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
			tc.setContext(c, tc.userID)
			c.Request = req

			userHandler.GetUserHandler(c)
			// assert
			tc.checkResponse(t, c, rec)
		})
	}
}
