package middleware

import (
	"go.opentelemetry.io/otel"

	liblog "github.com/kujilabo/cocotola-1.21/lib/log"
)

const (
	loggerKey = liblog.CoreControllerLoggerContextKey
)

var (
	tracer = otel.Tracer("github.com/kujilabo/cocotola-core/src/controller/gin/middleware")
)
