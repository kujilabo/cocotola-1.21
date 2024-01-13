package helper

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	rsliblog "github.com/kujilabo/redstart/lib/log"
	rsuserdomain "github.com/kujilabo/redstart/user/domain"

	liblog "github.com/kujilabo/cocotola-1.21/lib/log"
)

func HandleSecuredFunction(c *gin.Context, fn func(ctx context.Context, logger *slog.Logger, organizationID *rsuserdomain.OrganizationID, operatorID *rsuserdomain.AppUserID) error, errorHandle func(ctx context.Context, logger *slog.Logger, c *gin.Context, err error) bool) {
	ctx := c.Request.Context()
	authLogger := rsliblog.GetLoggerFromContext(ctx, liblog.AppAuthLoggerContextKey)

	organizationIDInt := c.GetInt("OrganizationID")
	if organizationIDInt == 0 {
		c.Status(http.StatusUnauthorized)
		return
	}

	organizationID, err := rsuserdomain.NewOrganizationID(organizationIDInt)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	appUserID := c.GetInt("AuthorizedUser")
	if appUserID == 0 {
		c.Status(http.StatusUnauthorized)
		return
	}

	operatorID, err := rsuserdomain.NewAppUserID(appUserID)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	authLogger.InfoContext(ctx, "", slog.Int("organization_id", organizationID.Int()), slog.Int("operator_id", operatorID.Int()))

	controllerLogger := rsliblog.GetLoggerFromContext(ctx, liblog.AppControllerLoggerContextKey)
	if err := fn(ctx, controllerLogger, organizationID, operatorID); err != nil {
		if handled := errorHandle(ctx, controllerLogger, c, err); !handled {
			c.Status(http.StatusInternalServerError)
		}
	}
}
