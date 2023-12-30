package handler_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/config"
	rsuserdomain "github.com/kujilabo/redstart/user/domain"

	handler "github.com/kujilabo/cocotola-1.21/cocotola-auth/src/controller/gin"
	handlermock "github.com/kujilabo/cocotola-1.21/cocotola-auth/src/controller/gin/mocks"
)

var (
	anythingOfContext = mock.MatchedBy(func(_ context.Context) bool { return true })
	corsConfig        cors.Config
	appConfig         *config.AppConfig
	authConfig        *config.AuthConfig
	debugConfig       *config.DebugConfig
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
	debugConfig = &config.DebugConfig{
		Gin:  false,
		Wait: false,
	}

	// signingKey := []byte(authConfig.SigningKey)
	// signingMethod := jwt.SigningMethodHS256
	// authTokenManager = auth.NewAuthTokenManager(signingKey, signingMethod, time.Duration(authConfig.AccessTokenTTLMin)*time.Minute, time.Duration(authConfig.RefreshTokenTTLHour)*time.Hour)
}

func initAuthRouter(t *testing.T, ctx context.Context, authentication handler.AuthenticationInterface) *gin.Engine {
	fn := handler.NewInitAuthRouterFunc(authentication)

	initPublicRouterFunc := []handler.InitRouterGroupFunc{fn}
	initPrivateRouterFunc := []handler.InitRouterGroupFunc{}

	router, err := handler.NewAppRouter(ctx, initPublicRouterFunc, initPrivateRouterFunc, corsConfig, appConfig, debugConfig)
	require.NoError(t, err)

	return router
}

func readBytes(t *testing.T, b *bytes.Buffer) []byte {
	respBytes, err := io.ReadAll(b)
	require.NoError(t, err)
	return respBytes
}

func parseJSON(t *testing.T, bytes []byte) interface{} {
	obj, err := oj.Parse(bytes)
	require.NoError(t, err)
	return obj
}

func parseExpr(t *testing.T, v string) jp.Expr {
	expr, err := jp.ParseString(v)
	require.NoError(t, err)
	return expr
}

func TestAuthHandler_GetUserInfo(t *testing.T) {
	ctx := context.Background()
	type conditions struct {
	}
	type inputs struct {
		authorizationHeader string
	}
	type outputs struct {
		statusCode     int
		appUserID      int
		organizationID int
		loginID        string
		username       string
	}
	type results struct {
	}
	tests := []struct {
		name       string
		conditions conditions
		inputs     inputs
		outputs    outputs
		results    results
	}{
		{
			name: "authorization header is not specified",
			outputs: outputs{
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name: "authorization header is invalid",
			inputs: inputs{
				authorizationHeader: "Bearer INVALID_TOKEN",
			},
			outputs: outputs{
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name: "authorization header is valid",
			inputs: inputs{
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
	appUserID, err := rsuserdomain.NewAppUserID(123)
	require.NoError(t, err)
	organizaionID, err := rsuserdomain.NewOrganizationID(456)
	require.NoError(t, err)
	appUserInfo := &rsuserdomain.AppUserModel{
		AppUserID:      appUserID,
		OrganizationID: organizaionID,
		LoginID:        "LOGIN_ID",
		Username:       "USERNAME",
	}
	authentication := new(handlermock.AuthenticationInterface)
	authentication.On("GetUserInfo", anythingOfContext, "VALID_TOKEN").Return(appUserInfo, nil)
	authentication.On("GetUserInfo", anythingOfContext, "INVALID_TOKEN").Return(nil, errors.New("INVALID"))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			r := initAuthRouter(t, ctx, authentication)
			w := httptest.NewRecorder()

			// when
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/v1/auth/userinfo", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", tt.inputs.authorizationHeader)
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
	ctx := context.Background()
	type conditions struct {
	}
	type inputs struct {
		requestBody string
	}
	type outputs struct {
		statusCode  int
		accessToken string
	}
	type results struct {
	}
	tests := []struct {
		name       string
		conditions conditions
		inputs     inputs
		outputs    outputs
		results    results
	}{
		{
			name: "requetyBody is empty",
			outputs: outputs{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "requetyBody is invalid",
			inputs: inputs{
				requestBody: `{"refreshToken": "INVALID_TOKEN"}`,
			},
			outputs: outputs{
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name: "requetyBody is valid",
			inputs: inputs{
				requestBody: `{"refreshToken": "VALID_TOKEN"}`,
			},
			outputs: outputs{
				statusCode:  http.StatusOK,
				accessToken: "ACCESS_TOKEN",
			},
		},
	}
	// given
	authentication := new(handlermock.AuthenticationInterface)
	authentication.On("RefreshToken", anythingOfContext, "VALID_TOKEN").Return("ACCESS_TOKEN", nil)
	authentication.On("RefreshToken", anythingOfContext, "INVALID_TOKEN").Return("", errors.New("INVALID"))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			r := initAuthRouter(t, ctx, authentication)
			w := httptest.NewRecorder()

			// when
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/v1/auth/refresh_token", bytes.NewBuffer([]byte(tt.inputs.requestBody)))
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
