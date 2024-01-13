package gateway

import (
	"context"
	"time"

	"gorm.io/gorm"

	rslibdomain "github.com/kujilabo/redstart/lib/domain"
	rsliberrors "github.com/kujilabo/redstart/lib/errors"
	rslibgateway "github.com/kujilabo/redstart/lib/gateway"

	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/service"
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

type RepositoryFactoryFunc func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error)
