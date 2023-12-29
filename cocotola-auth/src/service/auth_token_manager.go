package service

import (
	"context"

	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/domain"
	rsuserdomain "github.com/kujilabo/redstart/user/domain"
)

type AppUserInterface interface {
	AppUserID() *rsuserdomain.AppUserID
	OrganizationID() *rsuserdomain.OrganizationID
	LoginID() string
	Username() string
	// GetUserGroups() []domain.UserGroupModel
}

type OrganizationInterface interface {
	OrganizationID() *rsuserdomain.OrganizationID
	Name() string
}

type AuthTokenManager interface {
	CreateTokenSet(ctx context.Context, appUser AppUserInterface, organization OrganizationInterface) (*domain.AuthTokenSet, error)
	// RefreshToken(ctx context.Context, tokenString string) (string, error)
}
