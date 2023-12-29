//go:build small

package gateway_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/gateway"
	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/service"
	rsuserdomain "github.com/kujilabo/redstart/user/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type organization struct {
	organizationID *rsuserdomain.OrganizationID
	name           string
}

func (m *organization) OrganizationID() *rsuserdomain.OrganizationID {
	return m.organizationID
}
func (m *organization) Name() string {
	return m.name
}

type appUser struct {
	appUserID      *rsuserdomain.AppUserID
	organizationID *rsuserdomain.OrganizationID
	loginID        string
	username       string
}

func (m *appUser) AppUserID() *rsuserdomain.AppUserID {
	return m.appUserID
}
func (m *appUser) OrganizationID() *rsuserdomain.OrganizationID {
	return m.organizationID
}
func (m *appUser) Username() string {
	return m.username
}
func (m *appUser) LoginID() string {
	return m.loginID
}
func Test_authTokenManager_CreateTokenSet(t *testing.T) {
	organizationID, err := rsuserdomain.NewOrganizationID(123)
	require.NoError(t, err)
	appUserID, err := rsuserdomain.NewAppUserID(456)
	require.NoError(t, err)
	type fields struct {
		signingKey     []byte
		signingMethod  jwt.SigningMethod
		tokenTimeout   time.Duration
		refreshTimeout time.Duration
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
			name: "",
			fields: fields{
				signingKey:    []byte("&32KC^L;m'BuH+'ATNhv[qWM:3)2Lt2m"),
				signingMethod: jwt.SigningMethodHS256,
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
	}
	ctx := context.Background()
	for _, tt := range tests {
		authTokenManager := gateway.NewAuthTokenManager(tt.fields.signingKey, tt.fields.signingMethod, tt.fields.tokenTimeout, tt.fields.refreshTimeout)
		t.Run(tt.name, func(t *testing.T) {
			got, err := authTokenManager.CreateTokenSet(ctx, tt.args.appUser, tt.args.organization)
			if (err != nil) != tt.wantErr {
				t.Errorf("authTokenManager.CreateTokenSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotEmpty(t, got.AccessToken)
			assert.NotEmpty(t, got.RefreshToken)
		})
	}
}
