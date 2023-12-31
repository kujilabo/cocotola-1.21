package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	rsliblog "github.com/kujilabo/redstart/lib/log"
	rsuserdomain "github.com/kujilabo/redstart/user/domain"

	libapi "github.com/kujilabo/cocotola-1.21/lib/api"
	liblog "github.com/kujilabo/cocotola-1.21/lib/log"
)

type AuthenticationUsecaseInterface interface {
	GetUserInfo(ctx context.Context, bearerToken string) (*rsuserdomain.AppUserModel, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
}

type AuthHandler struct {
	authentication AuthenticationUsecaseInterface
}

func NewAuthHandler(authentication AuthenticationUsecaseInterface) *AuthHandler {
	return &AuthHandler{
		authentication: authentication,
	}
}

func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	logger := rsliblog.GetLoggerFromContext(ctx, liblog.AppControllerLoggerContextKey)
	logger.InfoContext(ctx, "GetUserInfo")

	authorization := c.GetHeader("Authorization")
	if !strings.HasPrefix(authorization, "Bearer ") {
		logger.WarnContext(ctx, "invalid header. Bearer not found")
		c.Status(http.StatusUnauthorized)
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

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	logger := rsliblog.GetLoggerFromContext(ctx, liblog.AppControllerLoggerContextKey)
	logger.InfoContext(ctx, "Authorize")
	refreshTokenParameter := RefreshTokenParameter{}
	if err := c.ShouldBindJSON(&refreshTokenParameter); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	accessToken, err := h.authentication.RefreshToken(ctx, refreshTokenParameter.RefreshToken)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		AccessToken: accessToken,
	})
}

func NewInitAuthRouterFunc(authentication AuthenticationUsecaseInterface) InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
		auth := parentRouterGroup.Group("auth")
		for _, m := range middleware {
			auth.Use(m)
		}

		authHandler := NewAuthHandler(authentication)
		auth.POST("refresh_token", authHandler.RefreshToken)
		auth.GET("userinfo", authHandler.GetUserInfo)
		return nil
	}
}
