package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	rsliblog "github.com/kujilabo/redstart/lib/log"

	libapi "github.com/kujilabo/cocotola-1.21/lib/api"
	liblog "github.com/kujilabo/cocotola-1.21/lib/log"

	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/domain"
)

type PasswordUsecaseInterface interface {
	Authenticate(ctx context.Context, loginID, password, organizationName string) (*domain.AuthTokenSet, error)
}
type PasswordAuthHandler struct {
	passwordUsecase PasswordUsecaseInterface
}

func NewPasswordAuthHandler(passwordUsecase PasswordUsecaseInterface) *PasswordAuthHandler {
	return &PasswordAuthHandler{
		passwordUsecase: passwordUsecase,
	}
}

func (h *PasswordAuthHandler) Authorize(c *gin.Context) {
	ctx := c.Request.Context()
	logger := rsliblog.GetLoggerFromContext(ctx, liblog.AppControllerLoggerContextKey)
	logger.Info("Authorize")

	passwordAuthParameter := libapi.PasswordAuthParameter{}
	if err := c.ShouldBindJSON(&passwordAuthParameter); err != nil {
		logger.InfoContext(ctx, fmt.Sprintf("invalid parameter. err: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest)})
		return
	}

	authResult, err := h.passwordUsecase.Authenticate(ctx, passwordAuthParameter.LoginID, passwordAuthParameter.Password, passwordAuthParameter.OrganizationName)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthenticated) {
			logger.InfoContext(ctx, fmt.Sprintf("invalid parameter. err: %v", err))
			c.JSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		logger.ErrorContext(ctx, fmt.Sprintf("passwordUsecase.Authenticate. err: %+v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, libapi.AuthResponse{
		AccessToken:  &authResult.AccessToken,
		RefreshToken: &authResult.RefreshToken,
	})
}

func NewInitPasswordRouterFunc(password PasswordUsecaseInterface) InitRouterGroupFunc {
	return func(parentRouterGroup gin.IRouter, middleware ...gin.HandlerFunc) error {
		auth := parentRouterGroup.Group("password")
		for _, m := range middleware {
			auth.Use(m)
		}

		passwordAuthHandler := NewPasswordAuthHandler(password)
		auth.POST("authenticate", passwordAuthHandler.Authorize)
		return nil
	}
}
