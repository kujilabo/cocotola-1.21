package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	rslibconfig "github.com/kujilabo/redstart/lib/config"
	rslibdomain "github.com/kujilabo/redstart/lib/domain"
	rsliberrors "github.com/kujilabo/redstart/lib/errors"
	rslibgateway "github.com/kujilabo/redstart/lib/gateway"
	rsliblog "github.com/kujilabo/redstart/lib/log"
	rssqls "github.com/kujilabo/redstart/sqls"
	rsusergateway "github.com/kujilabo/redstart/user/gateway"
	rsuserservice "github.com/kujilabo/redstart/user/service"

	libcontroller "github.com/kujilabo/cocotola-1.21/lib/controller/gin"
	liblog "github.com/kujilabo/cocotola-1.21/lib/log"
	"github.com/kujilabo/cocotola-1.21/proto"

	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/config"
	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/gateway"
	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/initialize"
	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/service"
	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/usecase"
)

const readHeaderTimeout = time.Duration(30) * time.Second

func getValue(values ...string) string {
	for _, v := range values {
		if len(v) != 0 {
			return v
		}
	}
	return ""
}

func main() {
	var _ = new(usecase.Authentication)
	ctx := context.Background()
	env := flag.String("env", "", "environment")
	flag.Parse()
	appEnv := getValue(*env, os.Getenv("APP_ENV"), "local")
	slog.InfoContext(ctx, fmt.Sprintf("env: %s", appEnv))

	rsliberrors.UseXerrorsErrorf()

	cfg, dialect, db, sqlDB, tp := Initialize(ctx, appEnv)
	defer sqlDB.Close()
	defer tp.ForceFlush(ctx) // flushes any pending spans

	ctx = liblog.InitLogger(ctx)
	logger := rsliblog.GetLoggerFromContext(ctx, rslibdomain.ContextKey(cfg.App.Name))

	rff := func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error) {
		return gateway.NewRepositoryFactory(ctx, dialect, cfg.DB.DriverName, db, time.UTC) // nolint:wrapcheck
	}
	rf, err := rff(ctx, db)
	if err != nil {
		panic(err)
	}
	rsrf, err := rsusergateway.NewRepositoryFactory(ctx, dialect, cfg.DB.DriverName, db, time.UTC)
	if err != nil {
		panic(err)
	}

	txManager, err := gateway.NewTransactionManager(db, rff)
	if err != nil {
		panic(err)
	}
	nonTxManager, err := gateway.NewNonTransactionManager(rf)
	if err != nil {
		panic(err)
	}

	logger.Info(fmt.Sprintf("%+v", proto.HelloRequest{}))

	initialize.InitApp1(ctx, txManager, nonTxManager, "cocotola", cfg.App.OwnerLoginID, cfg.App.OwnerPassword)

	gracefulShutdownTime2 := time.Duration(cfg.Shutdown.TimeSec2) * time.Second

	result := Run(ctx, cfg, txManager, nonTxManager, rsrf)

	time.Sleep(gracefulShutdownTime2)
	logger.InfoContext(ctx, "exited")
	os.Exit(result)
}

func Initialize(ctx context.Context, env string) (*config.Config, rslibgateway.DialectRDBMS, *gorm.DB, *sql.DB, *sdktrace.TracerProvider) {
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
	dialect, db, sqlDB, err := rslibconfig.InitDB(cfg.DB, rssqls.SQL)
	if err != nil {
		panic(err)
	}

	return cfg, dialect, db, sqlDB, tp
}

func Run(ctx context.Context, cfg *config.Config, txManager, nonTxManager service.TransactionManager, rsrf rsuserservice.RepositoryFactory) int {
	var eg *errgroup.Group
	eg, ctx = errgroup.WithContext(ctx)

	if !cfg.Debug.Gin {
		gin.SetMode(gin.ReleaseMode)
	}

	eg.Go(func() error {
		router := gin.New()
		if err := initialize.InitAppServer(ctx, router, cfg.CORS, cfg.Auth, cfg.Debug, cfg.App.Name, txManager, nonTxManager, rsrf); err != nil {
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
