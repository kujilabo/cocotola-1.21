//go:build small

package gateway_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/domain"
	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/gateway"
	gatewaymock "github.com/kujilabo/cocotola-1.21/cocotola-auth/src/gateway/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type errorType int

const (
	none errorType = iota
	other
	unauthenticated
)

func Test_GoogleAuthClient_NewGoogleAuthClient(t *testing.T) {
	t.Parallel()
	httpClient := new(gatewaymock.HTTPClient)
	c := gateway.NewGoogleAuthClient(
		httpClient, "CLIENT_ID", "CLIENT_SECRET", "REDIRECT_URI",
	)
	assert.Equal(t, "CLIENT_ID", c.ClientID)
	assert.Equal(t, "CLIENT_SECRET", c.ClientSecret)
	assert.Equal(t, "REDIRECT_URI", c.RedirectURI)
	assert.Equal(t, "authorization_code", c.GrantType)
	assert.Equal(t, httpClient, c.HTTPClient)
}

func Test_GoogleAuthClient_RetrieveAccessToken(t *testing.T) {
	ctx := context.Background()
	t.Parallel()
	type conditions struct {
		statusCode  int
		response    string
		doErrorType errorType
	}
	type args struct {
		code         string
		clientID     string
		clientSecret string
		redirectURI  string
		grantType    string
	}
	type outputs struct {
		accessToken   string
		refreshToken  string
		wantErrorType errorType
	}
	tests := []struct {
		name       string
		conditions conditions
		args       args
		outputs    outputs
	}{
		{
			name: "valid",
			conditions: conditions{
				statusCode: http.StatusOK,
				response:   `{"access_token":"ACCESS_TOKEN","refresh_token":"REFRESH_TOKEN"}`,
			},
			args: args{
				code:         "VALID_CODE",
				clientID:     "CLIENT_ID",
				clientSecret: "CLIENT_SECRET",
				redirectURI:  "REDIRECT_URI",
				grantType:    "GRANT_TYPE",
			},
			outputs: outputs{
				accessToken:  "ACCESS_TOKEN",
				refreshToken: "REFRESH_TOKEN",
			},
		},
		{
			name: "http.client.Do error",
			conditions: conditions{
				doErrorType: other,
			},
			args: args{
				code:         "VALID_CODE",
				clientID:     "CLIENT_ID",
				clientSecret: "CLIENT_SECRET",
				redirectURI:  "REDIRECT_URI",
				grantType:    "GRANT_TYPE",
			},
			outputs: outputs{
				accessToken:   "ACCESS_TOKEN",
				refreshToken:  "REFRESH_TOKEN",
				wantErrorType: other,
			},
		},
		{
			name: "status code is 401",
			conditions: conditions{
				statusCode: http.StatusUnauthorized,
			},
			args: args{
				code:         "VALID_CODE",
				clientID:     "CLIENT_ID",
				clientSecret: "CLIENT_SECRET",
				redirectURI:  "REDIRECT_URI",
				grantType:    "GRANT_TYPE",
			},
			outputs: outputs{
				accessToken:   "ACCESS_TOKEN",
				refreshToken:  "REFRESH_TOKEN",
				wantErrorType: unauthenticated,
			},
		},
		{
			name: "status code is 400",
			conditions: conditions{
				statusCode: http.StatusBadRequest,
			},
			args: args{
				code:         "VALID_CODE",
				clientID:     "CLIENT_ID",
				clientSecret: "CLIENT_SECRET",
				redirectURI:  "REDIRECT_URI",
				grantType:    "GRANT_TYPE",
			},
			outputs: outputs{
				accessToken:   "ACCESS_TOKEN",
				refreshToken:  "REFRESH_TOKEN",
				wantErrorType: other,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// given
			httpClient := new(gatewaymock.HTTPClient)
			if tt.conditions.doErrorType == none {
				httpClient.On("Do", mock.Anything).Return(&http.Response{
					StatusCode: tt.conditions.statusCode,
					Body:       io.NopCloser(strings.NewReader(tt.conditions.response)),
				}, nil)
			} else {
				httpClient.On("Do", mock.Anything).Return(nil, errors.New("ERROR"))
			}
			c := gateway.GoogleAuthClient{
				HTTPClient:   httpClient,
				ClientID:     tt.args.clientID,
				ClientSecret: tt.args.clientSecret,
				RedirectURI:  tt.args.redirectURI,
				GrantType:    tt.args.grantType,
			}

			// when
			tokenSet, err := c.RetrieveAccessToken(ctx, tt.args.code)
			if tt.outputs.wantErrorType == none {
				assert.NoError(t, err)
			} else if tt.outputs.wantErrorType == unauthenticated {
				assert.ErrorIs(t, err, domain.ErrUnauthenticated)
				return
			} else if tt.outputs.wantErrorType == other {
				assert.Error(t, err)
				assert.NotErrorIs(t, err, domain.ErrUnauthenticated)
				return
			}

			// then
			req, ok := httpClient.Calls[0].Arguments[0].(*http.Request)
			assert.True(t, ok)
			assert.Equal(t, http.MethodPost, req.Method)
			assert.Len(t, req.Header, 1)
			assert.Equal(t, "application/json", req.Header["Content-Type"][0])
			reqBytes, err := io.ReadAll(req.Body)
			require.NoError(t, err)
			assert.Equal(t, `{"client_id":"CLIENT_ID","client_secret":"CLIENT_SECRET","code":"VALID_CODE","grant_type":"GRANT_TYPE","redirect_uri":"REDIRECT_URI"}`, string(reqBytes))
			assert.Equal(t, "https", req.URL.Scheme)
			assert.Equal(t, "accounts.google.com", req.Host)
			assert.Equal(t, "/o/oauth2/token", req.URL.Path)

			assert.Equal(t, tt.outputs.accessToken, tokenSet.AccessToken)
			assert.Equal(t, tt.outputs.refreshToken, tokenSet.RefreshToken)
		})
	}
}
