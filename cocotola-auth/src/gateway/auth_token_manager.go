package gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/domain"
	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/service"
	liblog "github.com/kujilabo/cocotola-1.21/lib/log"

	rsliberrors "github.com/kujilabo/redstart/lib/errors"
	rsliblog "github.com/kujilabo/redstart/lib/log"
	rsuserdomain "github.com/kujilabo/redstart/user/domain"
)

type AppUserClaims struct {
	LoginID          string `json:"loginId"`
	AppUserID        int    `json:"appUserId"`
	Username         string `json:"username"`
	OrganizationID   int    `json:"organizationId"`
	OrganizationName string `json:"organizationName"`
	// Role             string `json:"role"`
	TokenType string `json:"tokenType"`
	jwt.StandardClaims
}

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

type authTokenManager struct {
	signingKey     []byte
	signingMethod  jwt.SigningMethod
	tokenTimeout   time.Duration
	refreshTimeout time.Duration
}

func NewAuthTokenManager(signingKey []byte, signingMethod jwt.SigningMethod, tokenTimeout, refreshTimeout time.Duration) service.AuthTokenManager {
	return &authTokenManager{
		signingKey:     signingKey,
		signingMethod:  signingMethod,
		tokenTimeout:   tokenTimeout,
		refreshTimeout: refreshTimeout,
	}
}

func (m *authTokenManager) CreateTokenSet(ctx context.Context, appUser service.AppUserInterface, organization service.OrganizationInterface) (*domain.AuthTokenSet, error) {
	accessToken, err := m.createJWT(ctx, appUser, organization, m.tokenTimeout, "access")
	if err != nil {
		return nil, err
	}

	refreshToken, err := m.createJWT(ctx, appUser, organization, m.refreshTimeout, "refresh")
	if err != nil {
		return nil, err
	}

	return &domain.AuthTokenSet{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (m *authTokenManager) createJWT(ctx context.Context, appUser service.AppUserInterface, organization service.OrganizationInterface, duration time.Duration, tokenType string) (string, error) {
	logger := rsliblog.GetLoggerFromContext(ctx, liblog.AuthGatewayLoggerContextKey)
	now := time.Now()
	claims := AppUserClaims{
		LoginID:          appUser.LoginID(),
		AppUserID:        appUser.AppUserID().Int(),
		Username:         appUser.Username(),
		OrganizationID:   organization.OrganizationID().Int(),
		OrganizationName: organization.Name(),
		TokenType:        tokenType,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(duration).Unix(),
		},
	}

	logger.DebugContext(ctx, fmt.Sprintf("claims: %+v", claims))

	token := jwt.NewWithClaims(m.signingMethod, claims)
	signed, err := token.SignedString(m.signingKey)
	if err != nil {
		return "", rsliberrors.Errorf(". err: %w", err)
	}

	return signed, nil
}

func (m *authTokenManager) RefreshToken(ctx context.Context, tokenString string) (string, error) {
	logger := rsliblog.GetLoggerFromContext(ctx, liblog.AuthGatewayLoggerContextKey)
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return m.signingKey, nil
	}

	currentToken, err := jwt.ParseWithClaims(tokenString, &AppUserClaims{}, keyFunc)
	if err != nil {
		logger.InfoContext(ctx, fmt.Sprintf("%v", err))
		return "", fmt.Errorf("jwt.ParseWithClaims. err: %w", domain.ErrUnauthenticated)
	}

	currentClaims, ok := currentToken.Claims.(*AppUserClaims)
	if !ok || !currentToken.Valid {
		return "", fmt.Errorf("invalid token. err: %w", domain.ErrUnauthenticated)
	}

	if currentClaims.TokenType != "refresh" {
		return "", fmt.Errorf("invalid token type. err: %w", domain.ErrUnauthenticated)
	}

	appUser := &appUser{
		loginID:  currentClaims.LoginID,
		username: currentClaims.Username,
	}
	organizationID, err := rsuserdomain.NewOrganizationID(currentClaims.OrganizationID)
	if err != nil {
		return "", err
	}

	organization := &organization{
		organizationID: organizationID,
		name:           currentClaims.OrganizationName,
	}

	accessToken, err := m.createJWT(ctx, appUser, organization, m.tokenTimeout, "access")
	if err != nil {
		return "", rsliberrors.Errorf("m.createJWT. err: %w", err)
	}

	return accessToken, nil
}
