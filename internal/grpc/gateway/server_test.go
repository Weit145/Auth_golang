package gateway_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Weit145/Auth_golang/internal/grpc/gateway"
	"github.com/Weit145/Auth_golang/internal/lib/logger/slogdiscard"
	"github.com/Weit145/Auth_golang/internal/service/current"
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

func TestRegistrationUser_Unit(t *testing.T) {
	tests := []struct {
		name          string
		token         string
		mockAsset     string
		mockRefresh   string
		mockError     error
		expectedErr   string
		serviceCalled bool
	}{
		{
			name:          "success",
			token:         "valid_token",
			mockAsset:     "test_access_token",
			mockRefresh:   "test_refresh_token",
			serviceCalled: true,
		},
		{
			name:          "empty token",
			token:         "",
			expectedErr:   "token is required",
			serviceCalled: false,
		},
		{
			name:          "Service error",
			token:         "valid_token",
			mockError:     errors.New("service confirm error"),
			expectedErr:   "failed to confirm user",
			serviceCalled: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockService := mocks.NewServiceAuth(t)

			if tc.serviceCalled {
				mockService.On("Confirm", mock.Anything, tc.token).
					Return(tc.mockAsset, tc.mockRefresh, tc.mockError).Once()
			}

			srv := newTestServer(t, mockService)

			req := &pb.TokenRequest{
				TokenPod: tc.token,
			}

			resp, err := srv.RegistrationUser(context.Background(), req)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tc.mockAsset, resp.AccessToken)
				require.NotNil(t, resp.Cookie)
				require.Equal(t, "refresh_token", resp.Cookie.Key)
				require.Equal(t, tc.mockRefresh, resp.Cookie.Value)
				require.True(t, resp.Cookie.Httponly)
				require.True(t, resp.Cookie.Secure)
				require.Equal(t, "lax", resp.Cookie.Samesite)
				require.Equal(t, int32(24), resp.Cookie.MaxAge)
			}

			if !tc.serviceCalled {
				mockService.AssertNotCalled(t, "Confirm", mock.Anything, mock.Anything)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestRefreshToken_Unit(t *testing.T) {
	tests := []struct {
		name          string
		refreshToken  string
		mockAsset     string
		mockError     error
		expectedErr   string
		serviceCalled bool
	}{
		{
			name:          "success",
			refreshToken:  "valid_refresh_token",
			mockAsset:     "new_access_token",
			serviceCalled: true,
		},
		{
			name:          "empty refresh token",
			refreshToken:  "",
			expectedErr:   "RefreshToken is required",
			serviceCalled: false,
		},
		{
			name:          "Service error",
			refreshToken:  "valid_refresh_token",
			mockAsset:     "",
			mockError:     errors.New("service refresh error"),
			expectedErr:   "failed to refresh token",
			serviceCalled: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockService := mocks.NewServiceAuth(t)

			if tc.serviceCalled {
				mockService.On("Refresh", mock.Anything, tc.refreshToken).
					Return(tc.mockAsset, tc.mockError).Once()
			}

			srv := newTestServer(t, mockService)

			req := &pb.CookieRequest{
				RefreshToken: tc.refreshToken,
			}

			resp, err := srv.RefreshToken(context.Background(), req)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tc.mockAsset, resp.AccessToken)
			}

			if !tc.serviceCalled {
				mockService.AssertNotCalled(t, "Refresh", mock.Anything, mock.Anything)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthenticate_Unit(t *testing.T) {
	tests := []struct {
		name          string
		login         string
		password      string
		mockAsset     string
		mockRefresh   string
		mockError     error
		expectedErr   string
		serviceCalled bool
	}{
		{
			name:          "success",
			login:         "test_login",
			password:      "test_password",
			mockAsset:     "test_access_token",
			mockRefresh:   "test_refresh_token",
			serviceCalled: true,
		},
		{
			name:          "empty login",
			login:         "",
			password:      "test_password",
			expectedErr:   "Login is required",
			serviceCalled: false,
		},
		{
			name:          "empty password",
			login:         "test_login",
			password:      "",
			expectedErr:   "Password is required",
			serviceCalled: false,
		},
		{
			name:          "Service error",
			login:         "test_login",
			password:      "test_password",
			mockAsset:     "",
			mockRefresh:   "",
			mockError:     errors.New("service login error"),
			expectedErr:   "failed to authenticate user",
			serviceCalled: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockService := mocks.NewServiceAuth(t)

			if tc.serviceCalled {
				mockService.On("LoginUser", mock.Anything, tc.login, tc.password).
					Return(tc.mockAsset, tc.mockRefresh, tc.mockError).Once()
			}

			srv := newTestServer(t, mockService)

			req := &pb.UserLoginRequest{
				Login:    tc.login,
				Password: tc.password,
			}

			resp, err := srv.Authenticate(context.Background(), req)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tc.mockAsset, resp.AccessToken)
				require.NotNil(t, resp.Cookie)
				require.Equal(t, "refresh_token", resp.Cookie.Key)
				require.Equal(t, tc.mockRefresh, resp.Cookie.Value)
				require.True(t, resp.Cookie.Httponly)
				require.True(t, resp.Cookie.Secure)
				require.Equal(t, "lax", resp.Cookie.Samesite)
				require.Equal(t, int32(24), resp.Cookie.MaxAge)
			}

			if !tc.serviceCalled {
				mockService.AssertNotCalled(t, "LoginUser", mock.Anything, mock.Anything, mock.Anything)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestCurrentUser_Unit(t *testing.T) {
	tests := []struct {
		name          string
		accessToken   string
		mockUser      *current.User
		mockError     error
		expectedErr   string
		serviceCalled bool
	}{
		{
			name:        "success",
			accessToken: "valid_access_token",
			mockUser: &current.User{
				Id:         1,
				Login:      "test_user",
				IsActive:   true,
				IsVerified: true,
				Role:       "user",
			},
			serviceCalled: true,
		},
		{
			name:          "empty access token",
			accessToken:   "",
			expectedErr:   "AssetToken is required",
			serviceCalled: false,
		},
		{
			name:          "Service error",
			accessToken:   "valid_access_token",
			mockUser:      nil,
			mockError:     errors.New("service current user error"),
			expectedErr:   "failed to get current user",
			serviceCalled: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockService := mocks.NewServiceAuth(t)

			if tc.serviceCalled {
				mockService.On("Current", mock.Anything, tc.accessToken).
					Return(tc.mockUser, tc.mockError).Once()
			}

			srv := newTestServer(t, mockService)

			req := &pb.UserCurrentRequest{
				AccessToken: tc.accessToken,
			}

			resp, err := srv.CurrentUser(context.Background(), req)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, int32(tc.mockUser.Id), resp.Id)
				require.Equal(t, tc.mockUser.Login, resp.Login)
				require.Equal(t, tc.mockUser.IsActive, resp.IsActive)
				require.Equal(t, tc.mockUser.IsVerified, resp.IsVerified)
				require.Equal(t, tc.mockUser.Role, resp.Role)
			}

			if !tc.serviceCalled {
				mockService.AssertNotCalled(t, "Current", mock.Anything, mock.Anything)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestLogOutUser_Unit(t *testing.T) {
	tests := []struct {
		name          string
		accessToken   string
		mockError     error
		expectedErr   string
		serviceCalled bool
	}{
		{
			name:          "success",
			accessToken:   "valid_access_token",
			serviceCalled: true,
		},
		{
			name:          "empty access token",
			accessToken:   "",
			expectedErr:   "AssetToken is required",
			serviceCalled: false,
		},
		{
			name:          "Service error",
			accessToken:   "valid_access_token",
			mockError:     errors.New("service logout error"),
			expectedErr:   "failed to Log out user",
			serviceCalled: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockService := mocks.NewServiceAuth(t)

			if tc.serviceCalled {
				mockService.On("LogOutUser", mock.Anything, tc.accessToken).
					Return(tc.mockError).Once()
			}

			srv := newTestServer(t, mockService)

			req := &pb.TokenRequest{
				TokenPod: tc.accessToken,
			}

			resp, err := srv.LogOutUser(context.Background(), req)

			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
			}

			if !tc.serviceCalled {
				mockService.AssertNotCalled(t, "LogOutUser", mock.Anything, mock.Anything)
			}

			mockService.AssertExpectations(t)
		})
	}
}
