package middleware

import (
	"log/slog"
	"strings"

	"github.com/gin-gonic/gin"

	rsliblog "github.com/kujilabo/redstart/lib/log"

	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/service"
	liblog "github.com/kujilabo/cocotola-1.21/lib/log"
)

func NewAuthMiddleware(cocotolaAuthClient service.CocotolaAuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ctx, span := tracer.Start(ctx, "authMiddleware")
		defer span.End()

		logger := rsliblog.GetLoggerFromContext(ctx, liblog.AppTraceLoggerContextKey)

		authorization := c.GetHeader("Authorization")
		if !strings.HasPrefix(authorization, "Bearer ") {
			logger.InfoContext(ctx, "invalid header. Bearer not found")
			return
		}

		bearerToken := authorization[len("Bearer "):]
		appUserInfo, err := cocotolaAuthClient.RetrieveUserInfo(ctx, bearerToken)
		if err != nil {
			logger.WarnContext(ctx, "getUserInfo")
			return
		}

		c.Set("AuthorizedUser", appUserInfo.AppUserID)
		c.Set("OrganizationID", appUserInfo.OrganizationID)

		logger.WarnContext(ctx, "authenticated", slog.Int("app_user_id", appUserInfo.AppUserID), slog.Int("organization_id", appUserInfo.OrganizationID))
	}
}
