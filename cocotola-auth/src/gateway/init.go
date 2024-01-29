package gateway

import (
	"go.opentelemetry.io/otel"

	liblog "github.com/kujilabo/cocotola-1.21/lib/log"
)

const (
	loggerKey = liblog.AuthGatewayLoggerContextKey
)

var (
	tracer = otel.Tracer("github.com/kujilabo/cocotola-1.21/cocotola-auth/gateway")
)
