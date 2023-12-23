package usecase

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt"

	rsuserdomain "github.com/kujilabo/redstart/user/domain"
	rsuserservice "github.com/kujilabo/redstart/user/service"

	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/domain"
	liblog "github.com/kujilabo/cocotola-1.21/lib/log"
	rsliblog "github.com/kujilabo/redstart/lib/log"
)

type AppUserClaims struct {
	LoginID          string `json:"loginId"`
	AppUserID        int    `json:"appUserId"`
	Username         string `json:"username"`
	OrganizationID   int    `json:"organizationId"`
	OrganizationName string `json:"organizationName"`
	TokenType        string `json:"tokenType"`
	jwt.StandardClaims
}

type Authentication struct {
	rf         rsuserservice.RepositoryFactory
	signingKey []byte
}

func NewAuthentication(rf rsuserservice.RepositoryFactory, signingKey []byte) *Authentication {
	return &Authentication{
		rf:         rf,
		signingKey: signingKey,
	}
}

func (u *Authentication) Authenticate(ctx context.Context, bearerToken string) (*rsuserdomain.AppUserModel, error) {
	logger := rsliblog.GetLoggerFromContext(ctx, liblog.AppUsecaseLoggerContextKey)

	token, err := jwt.ParseWithClaims(bearerToken, &AppUserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return u.signingKey, nil
	})
	if err != nil {
		logger.InfoContext(ctx, fmt.Sprintf("invalid token. err: %v", err))
		return nil, domain.ErrUnauthenticated
	}

	claims, ok := token.Claims.(*AppUserClaims)
	if !ok || !token.Valid {
		// logger.InfoContext(ctx, "invalid token")
		return nil, domain.ErrUnauthenticated
	}

	systemAdmin, err := rsuserservice.NewSystemAdmin(ctx, u.rf)
	if err != nil {
		return nil, err
	}

	organizationID, err := rsuserdomain.NewOrganizationID(claims.OrganizationID)
	if err != nil {
		return nil, err
	}

	appUserID, err := rsuserdomain.NewAppUserID(claims.AppUserID)
	if err != nil {
		return nil, err
	}

	// # TODO Check whether the token is registered in the Database

	appUserRepo := u.rf.NewAppUserRepository(ctx)
	systemOwner, err := appUserRepo.FindSystemOwnerByOrganizationID(ctx, systemAdmin, organizationID)
	if err != nil {
		return nil, err
	}

	appUser, err := systemOwner.FindAppUserByID(ctx, appUserID)
	if err != nil {
		return nil, err
	}

	return appUser.AppUserModel, nil
}
