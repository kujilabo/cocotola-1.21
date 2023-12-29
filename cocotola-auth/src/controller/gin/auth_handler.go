package handler

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	rsliblog "github.com/kujilabo/redstart/lib/log"

	libapi "github.com/kujilabo/cocotola-1.21/lib/api"
	liblog "github.com/kujilabo/cocotola-1.21/lib/log"

	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/usecase"
)

type AuthHandler interface {
	RefreshToken(c *gin.Context)
	GetUserInfo(c *gin.Context)
}

type authHandler struct {
	authentication *usecase.Authentication
}

func NewAuthHandler(authentication *usecase.Authentication) AuthHandler {
	return &authHandler{
		authentication: authentication,
	}
}

func (h *authHandler) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	logger := rsliblog.GetLoggerFromContext(ctx, liblog.AppControllerLoggerContextKey)
	logger.InfoContext(ctx, "GetUserInfo")

	authorization := c.GetHeader("Authorization")
	if !strings.HasPrefix(authorization, "Bearer ") {
		logger.WarnContext(ctx, "invalid header. Bearer not found")
		return
	}

	bearerToken := authorization[len("Bearer "):]
	appUserInfo, err := h.authentication.GetUserInfo(ctx, bearerToken)
	if err != nil {
		logger.WarnContext(ctx, "GetUserInfo", slog.Any("err", (err)))
		c.Status(http.StatusUnauthorized)
		return
	}
	c.JSON(http.StatusOK, libapi.AppUserInfoResponse{
		AppUserID:      appUserInfo.AppUserID.Int(),
		OrganizationID: appUserInfo.OrganizationID.Int(),
		LoginID:        appUserInfo.LoginID,
		Username:       appUserInfo.Username,
	})
	// TODO: check if the token is registered
}

func (h *authHandler) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	logger := rsliblog.GetLoggerFromContext(ctx, liblog.AppControllerLoggerContextKey)
	logger.InfoContext(ctx, "Authorize")
	refreshTokenParameter := RefreshTokenParameter{}
	if err := c.BindJSON(&refreshTokenParameter); err != nil {
		return
	}

	// token, err := h.authTokenManager.RefreshToken(ctx, refreshTokenParameter.RefreshToken)
	// if err != nil {
	// 	c.Status(http.StatusBadRequest)
	// 	return
	// }

	// logger.Info("Authorize OK")
	// c.JSON(http.StatusOK, AuthResponse{
	// 	AccessToken: token,
	// })
}
