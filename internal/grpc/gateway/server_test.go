package gateway_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Weit145/Auth_golang/internal/grpc/gateway"
	"github.com/Weit145/Auth_golang/internal/lib/logger/slogdiscard"
	"github.com/Weit145/Auth_golang/internal/service/mocks"
	pb "github.com/Weit145/proto-repo/auth"
)

func newTestServer(t *testing.T, svc *mocks.ServiceAuth) *gateway.Server {
	t.Helper()
	log := slogdiscard.NewDiscardLogger()
	return &gateway.Server{
		Service: svc,
		Log:     log,
	}
}

func TestCreateUser_Unit(t *testing.T) {
	tests := []struct {
		name          string
		login         string
		email         string
		password      string
		mockError     error
		expectedErr   string
		serviceCalled bool
	}{
		{
			name:          "success",
			login:         "test_user",
			email:         "test@example.com",
			password:      "password123",
			serviceCalled: true,
		},
		{
			name:        "empty login",
			email:       "test@example.com",
			password:    "password123",
			expectedErr: "login is required",
		},
		{
			name:        "empty email",
			login:       "test_user",
			password:    "password123",
			expectedErr: "email is required",
		},
		{
			name:        "empty password",
			login:       "test_user",
			email:       "test@example.com",
			expectedErr: "password is required",
		},
		{
			name:          "Service error",
			login:         "test_user",
			email:         "test@example.com",
			password:      "password123",
			mockError:     errors.New("db exploded"),
			expectedErr:   "failed to create user",
			serviceCalled: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockService := mocks.NewServiceAuth(t)

			if tc.serviceCalled {
				mockService.On("CreateUser", mock.Anything, tc.login, tc.email, tc.password).
					Return(tc.mockError).Once()
			}

			srv := newTestServer(t, mockService)

			req := &pb.UserCreateRequest{
				Login:    tc.login,
				Email:    tc.email,
				Password: tc.password,
			}

			resp, err := srv.CreateUser(context.Background(), req)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.True(t, resp.Success)
			}

			if !tc.serviceCalled {
				mockService.AssertNotCalled(t, "CreateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			}

			mockService.AssertExpectations(t)
		})
	}
}
