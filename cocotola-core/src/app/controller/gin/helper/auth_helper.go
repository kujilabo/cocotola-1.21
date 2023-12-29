package helper

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kujilabo/cocotola-1.21/lib/log"
	liblog "github.com/kujilabo/cocotola-1.21/lib/log"
	rsliblog "github.com/kujilabo/redstart/lib/log"
	"github.com/kujilabo/redstart/user/domain"
)

func HandleSecuredFunction(c *gin.Context, fn func(ctx context.Context, logger *slog.Logger, organizationID *domain.OrganizationID, operatorID *domain.AppUserID) error, errorHandle func(ctx context.Context, logger *slog.Logger, c *gin.Context, err error) bool) {
	ctx := c.Request.Context()
	authLogger := rsliblog.GetLoggerFromContext(ctx, liblog.AppAuthLoggerContextKey)

	organizationIDInt := c.GetInt("OrganizationID")
	if organizationIDInt == 0 {
		c.Status(http.StatusUnauthorized)
		return
	}

	organizationID, err := domain.NewOrganizationID(organizationIDInt)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	appUserID := c.GetInt("AuthorizedUser")
	if appUserID == 0 {
		c.Status(http.StatusUnauthorized)
		return
	}

	operatorID, err := domain.NewAppUserID(appUserID)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	authLogger.InfoContext(ctx, "", slog.Int("organization_id", organizationID.Int()), slog.Int("operator_id", operatorID.Int()))

	controllerLogger := rsliblog.GetLoggerFromContext(ctx, log.AppControllerLoggerContextKey)
	if err := fn(ctx, controllerLogger, organizationID, operatorID); err != nil {
		if handled := errorHandle(ctx, controllerLogger, c, err); !handled {
			c.Status(http.StatusInternalServerError)
		}
	}
}
