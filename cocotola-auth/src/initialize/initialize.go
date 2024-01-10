package initialize

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	rslibconfig "github.com/kujilabo/redstart/lib/config"
	rsliberrors "github.com/kujilabo/redstart/lib/errors"
	rslibgateway "github.com/kujilabo/redstart/lib/gateway"
	rsliblog "github.com/kujilabo/redstart/lib/log"
	rssqls "github.com/kujilabo/redstart/sqls"
	rsuserservice "github.com/kujilabo/redstart/user/service"

	libconfig "github.com/kujilabo/cocotola-1.21/lib/config"
	libcontroller "github.com/kujilabo/cocotola-1.21/lib/controller/gin"
	liblog "github.com/kujilabo/cocotola-1.21/lib/log"

	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/config"
	controller "github.com/kujilabo/cocotola-1.21/cocotola-auth/src/controller/gin"
	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/gateway"
	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/service"
	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/usecase"
)

const readHeaderTimeout = time.Duration(30) * time.Second

func Initialize(ctx context.Context, env string) (*config.Config, *gorm.DB, *sql.DB, *sdktrace.TracerProvider) {
	cfg, err := config.LoadConfig(env)
	if err != nil {
		panic(err)
	}

	// init log
	if err := rslibconfig.InitLog(cfg.Log); err != nil {
		panic(err)
	}

	// init tracer
	tp, err := rslibconfig.InitTracerProvider(ctx, cfg.App.Name, cfg.Trace)
	if err != nil {
		panic(err)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// init db
	db, sqlDB, err := rslibconfig.InitDB(cfg.DB, rssqls.SQL)
	if err != nil {
		panic(err)
	}

	return cfg, db, sqlDB, tp
}

func InitTransactionManager(db *gorm.DB, rff gateway.RepositoryFactoryFunc) service.TransactionManager {
	appTransactionManager, err := gateway.NewTransactionManager(db, rff)
	if err != nil {
		panic(err)
	}

	return appTransactionManager
}

func Run(ctx context.Context, cfg *config.Config, transactionManager service.TransactionManager, rsrf rsuserservice.RepositoryFactory) int {
	var eg *errgroup.Group
	eg, ctx = errgroup.WithContext(ctx)

	if !cfg.Debug.Gin {
		gin.SetMode(gin.ReleaseMode)
	}

	eg.Go(func() error {
		router := gin.New()
		if err := InitAppServer(ctx, router, cfg.CORS, cfg.Auth, cfg.Debug, cfg.App.Name, transactionManager, rsrf); err != nil {
			return err
		}
		return libcontroller.AppServerProcess(ctx, cfg.App.Name, router, cfg.App.HTTPPort, readHeaderTimeout, time.Duration(cfg.Shutdown.TimeSec1)*time.Second) // nolint:wrapcheck
	})
	eg.Go(func() error {
		return rslibgateway.MetricsServerProcess(ctx, cfg.App.MetricsPort, cfg.Shutdown.TimeSec1) // nolint:wrapcheck
	})
	eg.Go(func() error {
		return rslibgateway.SignalWatchProcess(ctx) // nolint:wrapcheck
	})
	eg.Go(func() error {
		<-ctx.Done()
		return ctx.Err() // nolint:wrapcheck
	})

	if err := eg.Wait(); err != nil {
		return 1
	}
	return 0
}

type systemOwnerByOrganizationName struct {
}

func (s systemOwnerByOrganizationName) Get(ctx context.Context, rf service.RepositoryFactory, organizationName string) (*rsuserservice.SystemOwner, error) {
	rsrf, err := rf.NewRedstartRepositoryFactory(ctx)
	if err != nil {
		return nil, err
	}
	systemAdmin, err := rsuserservice.NewSystemAdmin(ctx, rsrf)
	if err != nil {
		return nil, err
	}

	systemOwner, err := systemAdmin.FindSystemOwnerByOrganizationName(ctx, organizationName)
	if err != nil {
		return nil, err
	}

	return systemOwner, nil
}

func InitAppServer(ctx context.Context, parentRouterGroup gin.IRouter, corsConfig *rslibconfig.CORSConfig, authConfig *config.AuthConfig, debugConfig *libconfig.DebugConfig, appName string, transactionManager service.TransactionManager, rsrf rsuserservice.RepositoryFactory) error {
	// cors
	gincorsConfig := rslibconfig.InitCORS(corsConfig)
	httpClient := http.Client{
		Timeout:   time.Duration(authConfig.APITimeoutSec) * time.Second,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	googleAuthClient := gateway.NewGoogleAuthClient(&httpClient, authConfig.GoogleClientID, authConfig.GoogleClientSecret, authConfig.GoogleCallbackURL)

	signingKey := []byte(authConfig.SigningKey)
	signingMethod := jwt.SigningMethodHS256
	authTokenManager := gateway.NewAuthTokenManager(signingKey, signingMethod, time.Duration(authConfig.AccessTokenTTLMin)*time.Minute, time.Duration(authConfig.RefreshTokenTTLHour)*time.Hour)

	authenticationUsecase := usecase.NewAuthentication(transactionManager, authTokenManager, &systemOwnerByOrganizationName{})
	googleUserUsecase := usecase.NewGoogleUserUsecase(transactionManager, authTokenManager, googleAuthClient)

	privateRouterGroupFunc := []controller.InitRouterGroupFunc{}

	publicRouterGroupFunc := []controller.InitRouterGroupFunc{
		controller.NewInitTestRouterFunc(),
		controller.NewInitAuthRouterFunc(authenticationUsecase),
		controller.NewInitGoogleRouterFunc(googleUserUsecase),
	}

	if err := controller.InitRouter(ctx, parentRouterGroup, publicRouterGroupFunc, privateRouterGroupFunc, gincorsConfig, debugConfig, appName); err != nil {
		return err
	}

	return nil
}

func InitApp1(ctx context.Context, transactionManager service.TransactionManager, organizationName string, password string) {
	logger := rsliblog.GetLoggerFromContext(ctx, liblog.CoreMainLoggerContextKey)
	addOrganizationFunc := func(ctx context.Context, systemAdmin *rsuserservice.SystemAdmin) error {
		organization, err := systemAdmin.FindOrganizationByName(ctx, organizationName)
		if err == nil {
			logger.InfoContext(ctx, fmt.Sprintf("organization: %d", organization.OrganizationID().Int()))
			return nil
		} else if !errors.Is(err, rsuserservice.ErrOrganizationNotFound) {
			return rsliberrors.Errorf("failed to AddOrganization. err: %w", err)
		}

		firstOwnerAddParam, err := rsuserservice.NewAppUserAddParameter("cocotola-owner", "Owner(cocotola)", password, "", "", "", "")
		if err != nil {
			return rsliberrors.Errorf("NewFirstOwnerAddParameter. err: %w", err)
		}

		organizationAddParameter, err := rsuserservice.NewOrganizationAddParameter(organizationName, firstOwnerAddParam)
		if err != nil {
			return rsliberrors.Errorf("NewOrganizationAddParameter. err: %w", err)
		}

		organizationID, err := systemAdmin.AddOrganization(ctx, organizationAddParameter)
		if err != nil {
			return rsliberrors.Errorf("AddOrganization. err: %w", err)
		}

		logger.InfoContext(ctx, fmt.Sprintf("organizationID: %d", organizationID.Int()))
		return nil
	}

	if err := systemAdminAction(ctx, transactionManager, addOrganizationFunc); err != nil {
		panic(err)
	}
}

func systemAdminAction(ctx context.Context, transactionManager service.TransactionManager, fn func(context.Context, *rsuserservice.SystemAdmin) error) error {
	return transactionManager.Do(ctx, func(rf service.RepositoryFactory) error {
		rsrf, err := rf.NewRedstartRepositoryFactory(ctx)
		if err != nil {
			return rsliberrors.Errorf(". err: %w", err)
		}

		systemAdmin, err := rsuserservice.NewSystemAdmin(ctx, rsrf)
		if err != nil {
			return rsliberrors.Errorf(". err: %w", err)
		}

		return fn(ctx, systemAdmin)
	})
}
