package auth

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"gorm.io/gorm"

	rslibdomain "github.com/kujilabo/redstart/lib/domain"
	rsliberrors "github.com/kujilabo/redstart/lib/errors"
	rsliblog "github.com/kujilabo/redstart/lib/log"
	rsusergateway "github.com/kujilabo/redstart/user/gateway"

	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/gateway"
	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/initialize"
	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/service"
	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/usecase"
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
	var _ = new(usecase.Authentication)
	ctx := context.Background()
	env := flag.String("env", "", "environment")
	flag.Parse()
	appEnv := getValue(*env, os.Getenv("APP_ENV"), "local")
	slog.InfoContext(ctx, fmt.Sprintf("env: %s", appEnv))

	rsliberrors.UseXerrorsErrorf()

	cfg, db, sqlDB, tp := initialize.Initialize(ctx, appEnv)
	defer sqlDB.Close()
	defer tp.ForceFlush(ctx) // flushes any pending spans

	ctx = liblog.InitLogger(ctx)
	logger := rsliblog.GetLoggerFromContext(ctx, rslibdomain.ContextKey(cfg.App.Name))

	rff := func(ctx context.Context, db *gorm.DB) (service.RepositoryFactory, error) {
		return gateway.NewRepositoryFactory(ctx, cfg.DB.DriverName, db, time.UTC) // nolint:wrapcheck
	}
	rsrf, err := rsusergateway.NewRepositoryFactory(ctx, cfg.DB.DriverName, db, time.UTC)
	if err != nil {
		panic(err)
	}

	transactionManager := initialize.InitTransactionManager(db, rff)

	logger.Info(fmt.Sprintf("%+v", proto.HelloRequest{}))

	initialize.InitApp1(ctx, transactionManager, "cocotola", cfg.App.OwnerPassword)

	gracefulShutdownTime2 := time.Duration(cfg.Shutdown.TimeSec2) * time.Second

	result := initialize.Run(ctx, cfg, transactionManager, rsrf)

	time.Sleep(gracefulShutdownTime2)
	logger.InfoContext(ctx, "exited")
	os.Exit(result)
}
