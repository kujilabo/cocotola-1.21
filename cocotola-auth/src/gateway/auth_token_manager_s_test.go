//go:build small

package gateway_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/domain"
	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/gateway"
	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/service"
)

func Test_authTokenManager_CreateTokenSet(t *testing.T) {
	organizationID := organizationID(t, 123)
	appUserID := appUserID(t, 456)
	type fields struct {
		SigningKey     []byte
		SigningMethod  jwt.SigningMethod
		TokenTimeout   time.Duration
		RefreshTimeout time.Duration
	}
	type args struct {
		appUser      service.AppUserInterface
		organization service.OrganizationInterface
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				SigningKey:    []byte("&32KC^L;m'BuH+'ATNhv[qWM:3)2Lt2m"),
				SigningMethod: jwt.SigningMethodHS256,
			},
			args: args{
				appUser: &appUser{
					appUserID:      appUserID,
					organizationID: organizationID,
					loginID:        "LOGIN_ID",
					username:       "USERNAME",
				},
				organization: &organization{
					organizationID: organizationID,
					name:           "ORG_NAME",
				},
			},
			wantErr: false,
		},
		{
			name: "signingKey is not set",
			fields: fields{
				SigningKey:    []byte(""),
				SigningMethod: jwt.SigningMethodHS256,
			},
			args: args{
				appUser: &appUser{
					appUserID:      appUserID,
					organizationID: organizationID,
					loginID:        "LOGIN_ID",
					username:       "USERNAME",
				},
				organization: &organization{
					organizationID: organizationID,
					name:           "ORG_NAME",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		ctx := context.Background()
		m := &gateway.AuthTokenManager{
			SigningKey:     tt.fields.SigningKey,
			SigningMethod:  tt.fields.SigningMethod,
			TokenTimeout:   tt.fields.TokenTimeout,
			RefreshTimeout: tt.fields.RefreshTimeout,
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := m.CreateTokenSet(ctx, tt.args.appUser, tt.args.organization)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("authTokenManager.CreateTokenSet() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				return
			}
			assert.NotEmpty(t, got.AccessToken)
			assert.NotEmpty(t, got.RefreshToken)
		})
	}
}

func TestAuthTokenManager_GetUserInfo(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	organizationID := organizationID(t, 123)
	appUserID := appUserID(t, 456)
	appUser := &appUser{
		appUserID:      appUserID,
		organizationID: organizationID,
		loginID:        "LOGIN_ID",
		username:       "USERNAME",
	}
	organization := &organization{
		organizationID: organizationID,
		name:           "ORG_NAME",
	}

	type fields struct {
		SigningKey     []byte
		SigningMethod  jwt.SigningMethod
		TokenTimeout   time.Duration
		RefreshTimeout time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		want    *service.AppUserInfo
		wantErr error
	}{
		{
			name: "valid",
			fields: fields{
				SigningKey:    []byte("&32KC^L;m'BuH+'ATNhv[qWM:3)2Lt2m"),
				SigningMethod: jwt.SigningMethodHS256,
				TokenTimeout:  time.Second,
			},
			want: &service.AppUserInfo{
				LoginID:          "LOGIN_ID",
				Username:         "USERNAME",
				AppUserID:        456,
				OrganizationID:   123,
				OrganizationName: "ORG_NAME",
			},
			wantErr: nil,
		},
		{
			name: "expired",
			fields: fields{
				SigningKey:    []byte("&32KC^L;m'BuH+'ATNhv[qWM:3)2Lt2m"),
				SigningMethod: jwt.SigningMethodHS256,
				TokenTimeout:  -1 * time.Second,
			},
			wantErr: domain.ErrUnauthenticated,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := &gateway.AuthTokenManager{
				SigningKey:     tt.fields.SigningKey,
				SigningMethod:  tt.fields.SigningMethod,
				TokenTimeout:   tt.fields.TokenTimeout,
				RefreshTimeout: tt.fields.RefreshTimeout,
			}
			tokenSet, err := m.CreateTokenSet(ctx, appUser, organization)
			require.NoError(t, err)
			got, err := m.GetUserInfo(ctx, tokenSet.AccessToken)
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthTokenManager.GetUserInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
