package middleware

import "go.opentelemetry.io/otel"

var tracer = otel.Tracer("github.com/kujilabo/cocotola-core/src/app/controller/gin/middleware")
