package gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"

	rsliberrors "github.com/kujilabo/redstart/lib/errors"
	rsliblog "github.com/kujilabo/redstart/lib/log"
	rsuserdomain "github.com/kujilabo/redstart/user/domain"

	liblog "github.com/kujilabo/cocotola-1.21/lib/log"

	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/domain"
	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/service"
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

type AuthTokenManager struct {
	SigningKey     []byte
	SigningMethod  jwt.SigningMethod
	TokenTimeout   time.Duration
	RefreshTimeout time.Duration
}

func NewAuthTokenManager(signingKey []byte, signingMethod jwt.SigningMethod, tokenTimeout, refreshTimeout time.Duration) *AuthTokenManager {
	return &AuthTokenManager{
		SigningKey:     signingKey,
		SigningMethod:  signingMethod,
		TokenTimeout:   tokenTimeout,
		RefreshTimeout: refreshTimeout,
	}
}

func (m *AuthTokenManager) CreateTokenSet(ctx context.Context, appUser service.AppUserInterface, organization service.OrganizationInterface) (*domain.AuthTokenSet, error) {
	accessToken, err := m.createJWT(ctx, appUser, organization, m.TokenTimeout, "access")
	if err != nil {
		return nil, err
	}

	refreshToken, err := m.createJWT(ctx, appUser, organization, m.RefreshTimeout, "refresh")
	if err != nil {
		return nil, err
	}

	return &domain.AuthTokenSet{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (m *AuthTokenManager) createJWT(ctx context.Context, appUser service.AppUserInterface, organization service.OrganizationInterface, duration time.Duration, tokenType string) (string, error) {
	if len(m.SigningKey) == 0 {
		return "", rsliberrors.Errorf("m.SigningKey is not set")
	}

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

	token := jwt.NewWithClaims(m.SigningMethod, claims)
	signed, err := token.SignedString(m.SigningKey)
	if err != nil {
		return "", rsliberrors.Errorf(". err: %w", err)
	}

	return signed, nil
}

func (m *AuthTokenManager) GetUserInfo(ctx context.Context, tokenString string) (*service.AppUserInfo, error) {
	currentClaims, err := m.parseToken(ctx, tokenString)
	if err != nil {
		return nil, fmt.Errorf("parseToken(%s). err: %w", err.Error(), domain.ErrUnauthenticated)
	}

	return &service.AppUserInfo{
		LoginID:          currentClaims.LoginID,
		AppUserID:        currentClaims.AppUserID,
		Username:         currentClaims.Username,
		OrganizationID:   currentClaims.OrganizationID,
		OrganizationName: currentClaims.OrganizationName,
	}, nil
}

func (m *AuthTokenManager) parseToken(ctx context.Context, tokenString string) (*AppUserClaims, error) {
	logger := rsliblog.GetLoggerFromContext(ctx, liblog.AuthGatewayLoggerContextKey)
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return m.SigningKey, nil
	}

	currentToken, err := jwt.ParseWithClaims(tokenString, &AppUserClaims{}, keyFunc)
	if err != nil {
		logger.InfoContext(ctx, fmt.Sprintf("%v", err))
		// return nil, fmt.Errorf("jwt.ParseWithClaims. err: %w", domain.ErrUnauthenticated)
		return nil, err
	}
	if !currentToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	currentClaims, ok := currentToken.Claims.(*AppUserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	if err := currentClaims.Valid(); err != nil {
		return nil, err
	}

	return currentClaims, nil
}

func (m *AuthTokenManager) RefreshToken(ctx context.Context, tokenString string) (string, error) {
	currentClaims, err := m.parseToken(ctx, tokenString)
	if err != nil {
		return "", fmt.Errorf("parseToken(%s). err: %w", err.Error(), domain.ErrUnauthenticated)
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

	accessToken, err := m.createJWT(ctx, appUser, organization, m.TokenTimeout, "access")
	if err != nil {
		return "", rsliberrors.Errorf("m.createJWT. err: %w", err)
	}

	return accessToken, nil
}
