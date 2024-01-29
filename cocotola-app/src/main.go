package main

import (
	"context"
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"
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

	authgateway "github.com/kujilabo/cocotola-1.21/cocotola-auth/src/gateway"
	authinit "github.com/kujilabo/cocotola-1.21/cocotola-auth/src/initialize"
	authservice "github.com/kujilabo/cocotola-1.21/cocotola-auth/src/service"

	coregateway "github.com/kujilabo/cocotola-1.21/cocotola-core/src/gateway"
	coreinit "github.com/kujilabo/cocotola-1.21/cocotola-core/src/initialize"
	coreservice "github.com/kujilabo/cocotola-1.21/cocotola-core/src/service"
	coresqls "github.com/kujilabo/cocotola-1.21/cocotola-core/src/sqls"

	"github.com/kujilabo/cocotola-1.21/cocotola-app/src/config"
)

const (
	readHeaderTimeout = time.Duration(30) * time.Second

	loggerKey = liblog.AppMainLoggerContextKey
)

//go:embed web_dist
var web embed.FS

func getValue(values ...string) string {
	for _, v := range values {
		if len(v) != 0 {
			return v
		}
	}
	return ""
}

func main() {
	ctx := context.Background()
	env := flag.String("env", "", "environment")
	flag.Parse()
	appEnv := getValue(*env, os.Getenv("APP_ENV"), "local")
	slog.InfoContext(ctx, fmt.Sprintf("env: %s", appEnv))
	rsliberrors.UseXerrorsErrorf()

	cfg, dialect, db, sqlDB, tp := initialize(ctx, appEnv)
	defer sqlDB.Close()
	defer tp.ForceFlush(ctx) // flushes any pending spans

	ctx = liblog.InitLogger(ctx)
	ctx = rsliblog.WithLoggerName(ctx, loggerKey)
	logger := rsliblog.GetLoggerFromContext(ctx, rslibdomain.ContextKey(cfg.App.Name))

	authRFF := func(ctx context.Context, db *gorm.DB) (authservice.RepositoryFactory, error) {
		return authgateway.NewRepositoryFactory(ctx, dialect, cfg.DB.DriverName, db, time.UTC) // nolint:wrapcheck
	}
	authRF, err := authRFF(ctx, db)
	if err != nil {
		panic(err)
	}

	coreRFF := func(ctx context.Context, db *gorm.DB) (coreservice.RepositoryFactory, error) {
		return coregateway.NewRepositoryFactory(ctx, dialect, cfg.DB.DriverName, db, time.UTC) // nolint:wrapcheck
	}
	coreRF, err := coreRFF(ctx, db)
	if err != nil {
		panic(err)
	}

	rsrf, err := rsusergateway.NewRepositoryFactory(ctx, dialect, cfg.DB.DriverName, db, time.UTC)
	if err != nil {
		panic(err)
	}

	authTxManager, err := authgateway.NewTransactionManager(db, authRFF)
	if err != nil {
		panic(err)
	}

	authNonTxManager, err := authgateway.NewNonTransactionManager(authRF)
	if err != nil {
		panic(err)
	}

	coreTxManager, err := coregateway.NewTransactionManager(db, coreRFF)
	if err != nil {
		panic(err)
	}

	coreNonTxManager, err := coregateway.NewNonTransactionManager(coreRF)
	if err != nil {
		panic(err)
	}

	authinit.InitApp1(ctx, authTxManager, authNonTxManager, "cocotola", cfg.App.OwnerLoginID, cfg.App.OwnerPassword)

	gracefulShutdownTime2 := time.Duration(cfg.Shutdown.TimeSec2) * time.Second

	result := run(ctx, cfg, db, authTxManager, authNonTxManager, coreTxManager, coreNonTxManager, rsrf)

	time.Sleep(gracefulShutdownTime2)
	logger.InfoContext(ctx, "exited")
	os.Exit(result)
}

func initialize(ctx context.Context, env string) (*config.Config, rslibgateway.DialectRDBMS, *gorm.DB, *sql.DB, *sdktrace.TracerProvider) {
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
	dialect, db, sqlDB, err := rslibconfig.InitDB(cfg.DB, rssqls.SQL, coresqls.SQL)
	if err != nil {
		panic(err)
	}
	// if _, _, _, err := rslibconfig.InitDB(cfg.DB, coresqls.SQL); err != nil {
	// 	panic(err)
	// }

	return cfg, dialect, db, sqlDB, tp
}

func run(ctx context.Context, cfg *config.Config, db *gorm.DB, authTxManager authservice.TransactionManager, authNonTxManager authservice.TransactionManager, coreTxManager coreservice.TransactionManager, coreNonTxManager coreservice.TransactionManager, rsrf rsuserservice.RepositoryFactory) int {
	var eg *errgroup.Group
	eg, ctx = errgroup.WithContext(ctx)

	if !cfg.Debug.Gin {
		gin.SetMode(gin.ReleaseMode)
	}

	eg.Go(func() error {
		router := gin.New()

		viteStaticFS, err := fs.Sub(web, "web_dist")
		if err != nil {
			return err
		}

		router.NoRoute(func(c *gin.Context) {
			if strings.HasPrefix(c.Request.RequestURI, "/assets") {
				c.FileFromFS(c.Request.URL.Path, http.FS(viteStaticFS))
				return
			}
			if !strings.HasPrefix(c.Request.URL.Path, "/api") {
				c.FileFromFS("", http.FS(viteStaticFS))
				return
			}
		})

		api := router.Group("api")

		auth := api.Group("auth")
		if err := authinit.InitAppServer(ctx, auth, cfg.CORS, cfg.Auth, cfg.Debug, cfg.App.Name, authTxManager, authNonTxManager, rsrf); err != nil {
			return err
		}
		core := api.Group("core")
		if err := coreinit.InitAppServer(ctx, core, *cfg.AuthAPI, cfg.CORS, cfg.Debug, cfg.App.Name, db, coreTxManager, coreNonTxManager); err != nil {
			return err
		}

		return libcontroller.AppServerProcess(ctx, loggerKey, router, cfg.App.HTTPPort, readHeaderTimeout, time.Duration(cfg.Shutdown.TimeSec1)*time.Second) // nolint:wrapcheck
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
