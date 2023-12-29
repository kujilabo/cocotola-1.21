package service

import (
	"context"

	rsuserservice "github.com/kujilabo/redstart/user/service"
)

type RepositoryFactory interface {
	NewRedstartRepositoryFactory(ctx context.Context) (rsuserservice.RepositoryFactory, error)
}
