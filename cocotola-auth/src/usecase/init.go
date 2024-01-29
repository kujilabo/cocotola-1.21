package usecase

import (
	"context"

	rsuserservice "github.com/kujilabo/redstart/user/service"

	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/service"

	liblog "github.com/kujilabo/cocotola-1.21/lib/log"
)

const (
	loggerKey = liblog.AuthUsecaseLoggerContextKey
)

type SystemOwnerByOrganizationName interface {
	Get(ctx context.Context, rf service.RepositoryFactory, organizationName string) (*rsuserservice.SystemOwner, error)
}
