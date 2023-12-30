package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	liblog "github.com/kujilabo/cocotola-1.21/lib/log"

	rsliblog "github.com/kujilabo/redstart/lib/log"

	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/domain"
)

type GoogleAuthParameter struct {
	OrganizationName string `json:"organizationName"`
	Code             string `json:"code"`
}

type GoogleUserUsecaseInterface interface {
	RetrieveAccessToken(ctx context.Context, code string) (*domain.AuthTokenSet, error)

	RetrieveUserInfo(ctx context.Context, GoogleAuthResponse *domain.AuthTokenSet) (*domain.UserInfo, error)

	RegisterAppUser(ctx context.Context, googleUserInfo *domain.UserInfo, googleAuthResponse *domain.AuthTokenSet, organizationName string) (*domain.AuthTokenSet, error)
}

type GoogleUserHandler interface {
	Authorize(c *gin.Context)
}

type googleUserHandler struct {
	googleUserUsecase GoogleUserUsecaseInterface
}

func NewGoogleAuthHandler(googleUserUsecase GoogleUserUsecaseInterface) GoogleUserHandler {
	return &googleUserHandler{
		googleUserUsecase: googleUserUsecase,
	}
}

func (h *googleUserHandler) Authorize(c *gin.Context) {
	ctx := c.Request.Context()
	logger := rsliblog.GetLoggerFromContext(ctx, liblog.AppControllerLoggerContextKey)
	logger.Info("Authorize")

	googleAuthParameter := GoogleAuthParameter{}
	if err := c.BindJSON(&googleAuthParameter); err != nil {
		// logger.Warnf("invalid parameter. err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest)})
		return
	}

	// logger.Infof("RetrieveAccessToken. code: %s", googleAuthParameter)
	googleAuthResponse, err := h.googleUserUsecase.RetrieveAccessToken(ctx, googleAuthParameter.Code)
	if err != nil {
		// logger.Warnf("failed to RetrieveAccessToken. err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest)})
		return
	}

	// logger.Infof("RetrieveUserInfo. googleResponse: %+v", googleAuthResponse)
	userInfo, err := h.googleUserUsecase.RetrieveUserInfo(ctx, googleAuthResponse)
	if err != nil {
		// logger.Warnf("failed to RetrieveUserInfo. error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest)})
		return
	}

	logger.Info("RegisterAppUser")
	authResult, err := h.googleUserUsecase.RegisterAppUser(ctx, userInfo, googleAuthResponse, googleAuthParameter.OrganizationName)
	if err != nil {
		// logger.Warnf("failed to RegisterStudent. err: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": http.StatusText(http.StatusBadRequest)})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  authResult.AccessToken,
		RefreshToken: authResult.RefreshToken,
	})
}

func NewInitGoogleRouterFunc(googleUserUsecase GoogleUserUsecaseInterface) InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
		auth := parentRouterGroup.Group("google")
		for _, m := range middleware {
			auth.Use(m)
		}
		googleAuthHandler := NewGoogleAuthHandler(googleUserUsecase)
		auth.POST("authorize", googleAuthHandler.Authorize)
		return nil
	}
}
