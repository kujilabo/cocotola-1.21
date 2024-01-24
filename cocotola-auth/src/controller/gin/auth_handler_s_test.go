//go:build small

package handler_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	libconfig "github.com/kujilabo/cocotola-1.21/lib/config"

	rsuserdomain "github.com/kujilabo/redstart/user/domain"

	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/config"
	handler "github.com/kujilabo/cocotola-1.21/cocotola-auth/src/controller/gin"
	handlermock "github.com/kujilabo/cocotola-1.21/cocotola-auth/src/controller/gin/mocks"
)

var (
	anyOfCtx    = mock.MatchedBy(func(_ context.Context) bool { return true })
	corsConfig  cors.Config
	appConfig   *config.AppConfig
	authConfig  *config.AuthConfig
	debugConfig *libconfig.DebugConfig
	// authTokenManager  auth.AuthTokenManager
)

func init() {
	corsConfig = cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"*"},
		AllowHeaders:    []string{"*"},
	}
	appConfig = &config.AppConfig{
		Name:        "test",
		HTTPPort:    8080,
		MetricsPort: 8081,
	}
	authConfig = &config.AuthConfig{
		SigningKey:          "ah5T9Y9V2JPU74fhCtHQfDqLp3Zg8ZNc",
		AccessTokenTTLMin:   1,
		RefreshTokenTTLHour: 1,
	}
	debugConfig = &libconfig.DebugConfig{
		Gin:  false,
		Wait: false,
	}
}

func initAuthRouter(t *testing.T, ctx context.Context, authentication handler.AuthenticationUsecaseInterface) *gin.Engine {
	t.Helper()
	fn := handler.NewInitAuthRouterFunc(authentication)

	initPublicRouterFunc := []handler.InitRouterGroupFunc{fn}
	initPrivateRouterFunc := []handler.InitRouterGroupFunc{}

	router := gin.New()
	err := handler.InitRouter(ctx, router, initPublicRouterFunc, initPrivateRouterFunc, corsConfig, debugConfig, appConfig.Name)
	require.NoError(t, err)

	return router
}

func TestAuthHandler_GetUserInfo(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	type args struct {
		authorizationHeader string
	}
	type outputs struct {
		statusCode     int
		appUserID      int
		organizationID int
		loginID        string
		username       string
	}
	tests := []struct {
		name    string
		args    args
		outputs outputs
	}{
		{
			name: "authorization header is not specified",
			outputs: outputs{
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name: "authorization header is invalid",
			args: args{
				authorizationHeader: "Bearer INVALID_TOKEN",
			},
			outputs: outputs{
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name: "authorization header is valid",
			args: args{
				authorizationHeader: "Bearer VALID_TOKEN",
			},
			outputs: outputs{
				statusCode:     http.StatusOK,
				appUserID:      123,
				organizationID: 456,
				loginID:        "LOGIN_ID",
				username:       "USERNAME",
			},
		},
	}
	// given
	appUserID := appUserID(t, 123)
	organizaionID := organizationID(t, 456)
	appUserInfo := &rsuserdomain.AppUserModel{
		AppUserID:      appUserID,
		OrganizationID: organizaionID,
		LoginID:        "LOGIN_ID",
		Username:       "USERNAME",
	}
	authenticationUsecase := new(handlermock.AuthenticationUsecaseInterface)
	authenticationUsecase.On("GetUserInfo", anyOfCtx, "VALID_TOKEN").Return(appUserInfo, nil)
	authenticationUsecase.On("GetUserInfo", anyOfCtx, "INVALID_TOKEN").Return(nil, errors.New("INVALID"))

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// given
			r := initAuthRouter(t, ctx, authenticationUsecase)
			w := httptest.NewRecorder()

			// when
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/v1/auth/userinfo", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", tt.args.authorizationHeader)
			r.ServeHTTP(w, req)
			respBytes := readBytes(t, w.Body)

			// - check the status code
			assert.Equal(t, tt.outputs.statusCode, w.Code)

			if w.Code != http.StatusOK {
				assert.Len(t, respBytes, 0)
				return
			}

			jsonObj := parseJSON(t, respBytes)

			// - check the appUserId
			appUserIDExpr := parseExpr(t, "$.appUserId")
			appUserID := appUserIDExpr.Get(jsonObj)
			assert.Equal(t, int64(tt.outputs.appUserID), appUserID[0])

			// - check the organizationId
			organizationIDExpr := parseExpr(t, "$.organizationId")
			organizationID := organizationIDExpr.Get(jsonObj)
			assert.Equal(t, int64(tt.outputs.organizationID), organizationID[0])

			// - check the loginId
			loginIDExpr := parseExpr(t, "$.loginId")
			loginID := loginIDExpr.Get(jsonObj)
			assert.Equal(t, tt.outputs.loginID, loginID[0])

			// - check the username
			usernameExpr := parseExpr(t, "$.username")
			username := usernameExpr.Get(jsonObj)
			assert.Equal(t, tt.outputs.username, username[0])
		})
	}
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	type args struct {
		requestBody string
	}
	type outputs struct {
		statusCode  int
		accessToken string
	}
	tests := []struct {
		name    string
		args    args
		outputs outputs
	}{
		{
			name: "requetyBody is empty",
			outputs: outputs{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "requetyBody is invalid",
			args: args{
				requestBody: `{"refreshToken": "INVALID_TOKEN"}`,
			},
			outputs: outputs{
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name: "requetyBody is valid",
			args: args{
				requestBody: `{"refreshToken": "VALID_TOKEN"}`,
			},
			outputs: outputs{
				statusCode:  http.StatusOK,
				accessToken: "ACCESS_TOKEN",
			},
		},
	}
	// given
	authenticationUsecase := new(handlermock.AuthenticationUsecaseInterface)
	authenticationUsecase.On("RefreshToken", anyOfCtx, "VALID_TOKEN").Return("ACCESS_TOKEN", nil)
	authenticationUsecase.On("RefreshToken", anyOfCtx, "INVALID_TOKEN").Return("", errors.New("INVALID"))

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// given
			r := initAuthRouter(t, ctx, authenticationUsecase)
			w := httptest.NewRecorder()

			// when
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/v1/auth/refresh_token", bytes.NewBuffer([]byte(tt.args.requestBody)))
			require.NoError(t, err)
			r.ServeHTTP(w, req)

			respBytes := readBytes(t, w.Body)

			// - check the status code
			assert.Equal(t, tt.outputs.statusCode, w.Code)

			if w.Code != http.StatusOK {
				assert.Len(t, respBytes, 0)
				return
			}

			jsonObj := parseJSON(t, respBytes)

			// - check the organizationId
			accessTokenExpr := parseExpr(t, "$.accessToken")
			accessToken := accessTokenExpr.Get(jsonObj)
			assert.Equal(t, tt.outputs.accessToken, accessToken[0])
		})
	}
}
