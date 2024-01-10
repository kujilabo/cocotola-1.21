package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	liblog "github.com/kujilabo/cocotola-1.21/lib/log"

	rsliblog "github.com/kujilabo/redstart/lib/log"

	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/domain"
)

type googleAuthParameter struct {
	OrganizationName string `json:"organizationName"`
	Code             string `json:"code"`
}

type GoogleUserUsecaseInterface interface {
	Authorize(ctx context.Context, code, organizationName string) (*domain.AuthTokenSet, error)
}

type GoogleUserHandler struct {
	googleUserUsecase GoogleUserUsecaseInterface
}

func NewGoogleAuthHandler(googleUserUsecase GoogleUserUsecaseInterface) *GoogleUserHandler {
	return &GoogleUserHandler{
		googleUserUsecase: googleUserUsecase,
	}
}

func (h *GoogleUserHandler) Authorize(c *gin.Context) {
	ctx := c.Request.Context()
	logger := rsliblog.GetLoggerFromContext(ctx, liblog.AppControllerLoggerContextKey)
	logger.Info("Authorize")

	googleAuthParameter := googleAuthParameter{}
	if err := c.ShouldBindJSON(&googleAuthParameter); err != nil {
		// logger.Warnf("invalid parameter. err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest)})
		return
	}

	// // logger.Infof("RetrieveAccessToken. code: %s", googleAuthParameter)
	// googleAuthResponse, err := h.googleUserUsecase.RetrieveAccessToken(ctx, googleAuthParameter.Code)
	// if err != nil {
	// 	// logger.Warnf("failed to RetrieveAccessToken. err: %v", err)
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest)})
	// 	return
	// }

	// // logger.Infof("RetrieveUserInfo. googleResponse: %+v", googleAuthResponse)
	// userInfo, err := h.googleUserUsecase.RetrieveUserInfo(ctx, googleAuthResponse)
	// if err != nil {
	// 	// logger.Warnf("failed to RetrieveUserInfo. error: %v", err)
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": http.StatusText(http.StatusBadRequest)})
	// 	return
	// }

	// logger.Info("RegisterAppUser")
	// authResult, err := h.googleUserUsecase.RegisterAppUser(ctx, userInfo, googleAuthResponse, googleAuthParameter.OrganizationName)
	// if err != nil {
	// 	// logger.Warnf("failed to RegisterStudent. err: %+v", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": http.StatusText(http.StatusBadRequest)})
	// 	return
	// }

	if googleAuthParameter.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "code is empty"})
		return
	}
	if googleAuthParameter.OrganizationName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "organizationName is empty"})
		return
	}

	authResult, err := h.googleUserUsecase.Authorize(ctx, googleAuthParameter.Code, googleAuthParameter.OrganizationName)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthenticated) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		logger.ErrorContext(ctx, fmt.Sprintf("failed to RegisterStudent. err: %+v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  &authResult.AccessToken,
		RefreshToken: &authResult.RefreshToken,
	})
}

func NewInitGoogleRouterFunc(googleUser GoogleUserUsecaseInterface) InitRouterGroupFunc {
	return func(parentRouterGroup gin.IRouter, middleware ...gin.HandlerFunc) error {
		auth := parentRouterGroup.Group("google")
		for _, m := range middleware {
			auth.Use(m)
		}

		googleAuthHandler := NewGoogleAuthHandler(googleUser)
		auth.POST("authorize", googleAuthHandler.Authorize)
		return nil
	}
}
