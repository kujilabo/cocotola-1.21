package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/config"

	handler "github.com/kujilabo/cocotola-1.21/cocotola-auth/src/controller/gin"
	handlermock "github.com/kujilabo/cocotola-1.21/cocotola-auth/src/controller/gin/mocks"
)

var (
	// anythingOfContext = mock.MatchedBy(func(_ context.Context) bool { return true })
	corsConfig  cors.Config
	appConfig   *config.AppConfig
	authConfig  *config.AuthConfig
	debugConfig *config.DebugConfig
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

func TestAuthHandler_GetUserInfo(t *testing.T) {
	ctx := context.Background()
	authentication := new(handlermock.AuthenticationInterface)
	type conditions struct {
	}
	type inputs struct {
		bearerToken string
	}
	type outputs struct {
		statusCode int
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
			name: "BearerToken is not specified",
			outputs: outputs{
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name: "BearerToken is invalid",
			inputs: inputs{
				bearerToken: "INVALID_TOKEN",
			},
			outputs: outputs{
				statusCode: http.StatusUnauthorized,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := initAuthRouter(t, ctx, authentication)
			w := httptest.NewRecorder()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/v1/auth/userinfo", nil)
			require.NoError(t, err)
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.outputs.statusCode, w.Code)
		})
	}
}
