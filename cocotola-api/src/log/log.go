package log

import (
	"context"
	"log/slog"

	rslibdomain "github.com/kujilabo/redstart/lib/domain"
	rsliblog "github.com/kujilabo/redstart/lib/log"
)

const (
	LibGatewayLoggerContextKey    rslibdomain.ContextKey = "lib_gateway"
	AppGORMLoggerContextKey       rslibdomain.ContextKey = "app_gorm"
	AppServiceLoggerContextKey    rslibdomain.ContextKey = "app_service"
	AppGatewayLoggerContextKey    rslibdomain.ContextKey = "app_gateway"
	AppControllerLoggerContextKey rslibdomain.ContextKey = "app_controller"
	AppGinLoggerContextKey        rslibdomain.ContextKey = "app_gin"
	AppTraceLoggerContextKey      rslibdomain.ContextKey = "app_trace"
	AppAuthLoggerContextKey       rslibdomain.ContextKey = "app_auth"
)

var (
	LoggerKeys = []rslibdomain.ContextKey{
		LibGatewayLoggerContextKey,
		AppGORMLoggerContextKey,
		AppServiceLoggerContextKey,
		AppGatewayLoggerContextKey,
		AppControllerLoggerContextKey,
		AppGinLoggerContextKey,
		AppTraceLoggerContextKey,
		AppAuthLoggerContextKey,
	}
)

func InitLogger(ctx context.Context) context.Context {
	for _, key := range LoggerKeys {
		if _, ok := rsliblog.Loggers[key]; !ok {
			rsliblog.Loggers[key] = slog.New(rsliblog.LogHandlers[rsliblog.DefaultLogLevel])
		}
		ctx = context.WithValue(ctx, key, rsliblog.Loggers[key])
	}
	return ctx
}
