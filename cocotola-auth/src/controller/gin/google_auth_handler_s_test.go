//go:build small

package handler_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	handler "github.com/kujilabo/cocotola-1.21/cocotola-auth/src/controller/gin"
	handlermock "github.com/kujilabo/cocotola-1.21/cocotola-auth/src/controller/gin/mocks"
	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/domain"
)

func initGoogleRouter(t *testing.T, ctx context.Context, googleUser handler.GoogleUserUsecaseInterface) *gin.Engine {
	t.Helper()
	fn := handler.NewInitGoogleRouterFunc(googleUser)

	initPublicRouterFunc := []handler.InitRouterGroupFunc{fn}
	initPrivateRouterFunc := []handler.InitRouterGroupFunc{}

	router, err := handler.NewAppRouter(ctx, initPublicRouterFunc, initPrivateRouterFunc, corsConfig, appConfig, debugConfig)
	require.NoError(t, err)

	return router
}

func TestGoogleAuthHandler_Authorize(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	type args struct {
		requestBody string
	}
	type outputs struct {
		statusCode   int
		accessToken  string
		refreshToken string
	}
	tests := []struct {
		name    string
		args    args
		outputs outputs
	}{
		{
			name: "request body is empty",
			outputs: outputs{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "request body is invalid",
			args: args{
				requestBody: "[]",
			},
			outputs: outputs{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "code is invalid(error)",
			args: args{
				requestBody: `{"organizationName": "ORGANIZATION_NAME", "code": "ERROR_CODE"}`,
			},
			outputs: outputs{
				statusCode: http.StatusInternalServerError,
			},
		},
		{
			name: "code is invalid(unauthenticated)",
			args: args{
				requestBody: `{"organizationName": "ORGANIZATION_NAME", "code": "UNAUTHENTICATED_CODE"}`,
			},
			outputs: outputs{
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name: "code is valid",
			args: args{
				requestBody: `{"organizationName": "ORGANIZATION_NAME", "code": "VALID_CODE"}`,
			},
			outputs: outputs{
				statusCode:   http.StatusOK,
				accessToken:  "ACCESS_TOKEN",
				refreshToken: "REFRESH_TOKEN",
			},
		},
	}
	// given
	authToksenSet := &domain.AuthTokenSet{
		AccessToken:  "ACCESS_TOKEN",
		RefreshToken: "REFRESH_TOKEN",
	}
	googleUserUsecase := new(handlermock.GoogleUserUsecaseInterface)
	googleUserUsecase.On("Authorize", anythingOfContext, "VALID_CODE", "ORGANIZATION_NAME").Return(authToksenSet, nil)
	googleUserUsecase.On("Authorize", anythingOfContext, "ERROR_CODE", "ORGANIZATION_NAME").Return(nil, errors.New("INVALID"))
	googleUserUsecase.On("Authorize", anythingOfContext, "UNAUTHENTICATED_CODE", "ORGANIZATION_NAME").Return(nil, domain.ErrUnauthenticated)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// given
			r := initGoogleRouter(t, ctx, googleUserUsecase)
			w := httptest.NewRecorder()

			// when
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/v1/google/authorize", bytes.NewBuffer([]byte(tt.args.requestBody)))
			require.NoError(t, err)
			r.ServeHTTP(w, req)
			respBytes := readBytes(t, w.Body)

			// - check the status code
			assert.Equal(t, tt.outputs.statusCode, w.Code)

			jsonObj := parseJSON(t, respBytes)

			if w.Code != http.StatusOK {
				// - check the message
				messageExpr := parseExpr(t, "$.message")
				message := messageExpr.Get(jsonObj)
				assert.Equal(t, http.StatusText(w.Code), message[0])
				return
			}

			// - check the accessToken
			accessTokenExpr := parseExpr(t, "$.accessToken")
			accessToken := accessTokenExpr.Get(jsonObj)
			assert.Equal(t, tt.outputs.accessToken, accessToken[0])

			// - check the refreshToken
			refreshTokenExpr := parseExpr(t, "$.refreshToken")
			refreshToken := refreshTokenExpr.Get(jsonObj)
			assert.Equal(t, tt.outputs.refreshToken, refreshToken[0])
		})
	}
}
