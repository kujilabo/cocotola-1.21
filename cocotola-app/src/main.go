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
	rsuserservice "github.com/kujilabo/redstart/user/service"

	libcontroller "github.com/kujilabo/cocotola-1.21/lib/controller/gin"
	liblog "github.com/kujilabo/cocotola-1.21/lib/log"
	rsusergateway "github.com/kujilabo/redstart/user/gateway"

	authgateway "github.com/kujilabo/cocotola-1.21/cocotola-auth/src/gateway"
	authinit "github.com/kujilabo/cocotola-1.21/cocotola-auth/src/initialize"
	authservice "github.com/kujilabo/cocotola-1.21/cocotola-auth/src/service"

	"github.com/kujilabo/cocotola-1.21/cocotola-app/src/config"
)

const readHeaderTimeout = time.Duration(30) * time.Second

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

	cfg, db, sqlDB, tp := initialize(ctx, appEnv)
	defer sqlDB.Close()
	defer tp.ForceFlush(ctx) // flushes any pending spans

	ctx = liblog.InitLogger(ctx)
	logger := rsliblog.GetLoggerFromContext(ctx, rslibdomain.ContextKey(cfg.App.Name))
	rff := func(ctx context.Context, db *gorm.DB) (authservice.RepositoryFactory, error) {
		return authgateway.NewRepositoryFactory(ctx, cfg.DB.DriverName, db, time.UTC) // nolint:wrapcheck
	}
	rsrf, err := rsusergateway.NewRepositoryFactory(ctx, cfg.DB.DriverName, db, time.UTC)
	if err != nil {
		panic(err)
	}

	authTransactionManager := authinit.InitTransactionManager(db, rff)

	gracefulShutdownTime2 := time.Duration(cfg.Shutdown.TimeSec2) * time.Second

	result := run(ctx, cfg, authTransactionManager, rsrf)

	time.Sleep(gracefulShutdownTime2)
	logger.InfoContext(ctx, "exited")
	os.Exit(result)
}

func initialize(ctx context.Context, env string) (*config.Config, *gorm.DB, *sql.DB, *sdktrace.TracerProvider) {
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

func run(ctx context.Context, cfg *config.Config, authTransactionManager authservice.TransactionManager, rsrf rsuserservice.RepositoryFactory) int {
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
		if err := authinit.InitAppServer(ctx, auth, cfg.CORS, cfg.Auth, cfg.Debug, cfg.App.Name, authTransactionManager, rsrf); err != nil {
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