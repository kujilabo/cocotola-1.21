package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
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

	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/app/config"
	controller "github.com/kujilabo/cocotola-1.21/cocotola-core/src/app/controller/gin"
	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/app/domain"
	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/app/gateway"
	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/app/service"
	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/sqls"
	liblog "github.com/kujilabo/cocotola-1.21/lib/log"
	"github.com/kujilabo/cocotola-1.21/proto"
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

	rff := func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error) {
		return gateway.NewRepositoryFactory(ctx, cfg.DB.DriverName, db, time.UTC) // nolint:wrapcheck
	}

	appTransactionManager := initTransactionManager(db, rff)

	logger.Info(fmt.Sprintf("%+v", proto.HelloRequest{}))

	logger.Info("")
	logger.Info(domain.Lang2EN.String())
	logger.Info("Hello")
	logger.Info("Hello")
	service.A()

	gracefulShutdownTime2 := time.Duration(cfg.Shutdown.TimeSec2) * time.Second

	result := run(ctx, cfg, appTransactionManager)

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
	db, sqlDB, err := rslibconfig.InitDB(cfg.DB, sqls.SQL)
	if err != nil {
		panic(err)
	}

	return cfg, db, sqlDB, tp
}

func initTransactionManager(db *gorm.DB, rff gateway.RepositoryFactoryFunc) service.TransactionManager {
	appTransactionManager, err := gateway.NewTransactionManager(db, rff)
	if err != nil {
		panic(err)
	}

	return appTransactionManager
}

func run(ctx context.Context, cfg *config.Config, transactionManager service.TransactionManager) int {
	var eg *errgroup.Group
	eg, ctx = errgroup.WithContext(ctx)

	if !cfg.Debug.GinMode {
		gin.SetMode(gin.ReleaseMode)
	}

	eg.Go(func() error {
		return appServer(ctx, cfg) // nolint:wrapcheck
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

func appServer(ctx context.Context, cfg *config.Config) error {
	logger := rsliblog.GetLoggerFromContext(ctx, rslibdomain.ContextKey(cfg.App.Name))
	// // cors
	corsConfig := rslibconfig.InitCORS(cfg.CORS)
	// logrus.Infof("cors: %+v", corsConfig)

	// if err := corsConfig.Validate(); err != nil {
	// 	return liberrors.Errorf("corsConfig.Validate. err: %w", err)
	// }

	// studyMonitor := service.NewStudyMonitor()
	// studyStatUpdater := studyStatUpdater{
	// 	systemOwnerModel: systemOwnerModel,
	// 	appTransaction:   appTransaction,
	// }
	// if err := studyMonitor.Attach(&studyStatUpdater); err != nil {
	// 	return liberrors.Errorf(". err: %w", err)
	// }

	// privateRouterGroupFunc := []controller.InitRouterGroupFunc{
	// 	controller.NewInitWorkbookRouterFunc(studentUsecaseWorkbook),
	// 	controller.NewInitProblemRouterFunc(studentUsecaseProblem, newIteratorFunc),
	// 	controller.NewInitStudyRouterFunc(studentUseCaseStudy),
	// 	controller.NewInitAudioRouterFunc(studentUsecaseAudio),
	// 	controller.NewInitStatRouterFunc(studentUsecaseStat),
	// }

	publicRouterGroupFunc := []controller.InitRouterGroupFunc{
		controller.NewInitTestRouterFunc(),
	}
	router, err := controller.NewAppRouter(ctx,
		publicRouterGroupFunc,
		//privateRouterGroupFunc, pluginRouterGroupFunc, authTokenManager,
		corsConfig, cfg.App,
		//cfg.Auth,
		cfg.Debug)
	if err != nil {
		panic(err)
	}

	// if cfg.Swagger.Enabled {
	// 	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 	docs.SwaggerInfo.Title = cfg.App.Name
	// 	docs.SwaggerInfo.Version = "1.0"
	// 	docs.SwaggerInfo.Host = cfg.Swagger.Host
	// 	docs.SwaggerInfo.Schemes = []string{cfg.Swagger.Schema}
	// }

	httpServer := http.Server{
		Addr:              ":" + strconv.Itoa(cfg.App.HTTPPort),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	logger.InfoContext(ctx, fmt.Sprintf("http server listening at %v", httpServer.Addr))

	errCh := make(chan error)
	go func() {
		defer close(errCh)
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.InfoContext(ctx, fmt.Sprintf("failed to ListenAndServe. err: %v", err))
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		gracefulShutdownTime1 := time.Duration(cfg.Shutdown.TimeSec1) * time.Second
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), gracefulShutdownTime1)
		defer shutdownCancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.InfoContext(ctx, fmt.Sprintf("Server forced to shutdown. err: %v", err))
			return rsliberrors.Errorf(". err: %w", err)
		}
		return nil
	case err := <-errCh:
		return rsliberrors.Errorf(". err: %w", err)
	}
}
