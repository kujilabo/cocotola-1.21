package gateway

import (
	"context"
	"time"

	"gorm.io/gorm"

	rslibdomain "github.com/kujilabo/redstart/lib/domain"
	rsliberrors "github.com/kujilabo/redstart/lib/errors"
	rslibgateway "github.com/kujilabo/redstart/lib/gateway"
	rsusergateway "github.com/kujilabo/redstart/user/gateway"
	rsuserservice "github.com/kujilabo/redstart/user/service"

	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/service"
)

type RepositoryFactory struct {
	dialect    rslibgateway.DialectRDBMS
	driverName string
	db         *gorm.DB
	location   *time.Location
}

func NewRepositoryFactory(ctx context.Context, dialect rslibgateway.DialectRDBMS, driverName string, db *gorm.DB, location *time.Location) (*RepositoryFactory, error) {
	if db == nil {
		return nil, rsliberrors.Errorf("db is nil. err: %w", rslibdomain.ErrInvalidArgument)
	}

	return &RepositoryFactory{
		dialect:    dialect,
		driverName: driverName,
		db:         db,
		location:   location,
	}, nil
}

func (f *RepositoryFactory) NewRedstartRepositoryFactory(ctx context.Context) (rsuserservice.RepositoryFactory, error) {
	return rsusergateway.NewRepositoryFactory(ctx, f.dialect, f.driverName, f.db, f.location)
}

func (f *RepositoryFactory) NewStateRepository(ctx context.Context) (service.StateRepository, error) {
	return NewStateRepository(ctx)
}

type RepositoryFactoryFunc func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error)
