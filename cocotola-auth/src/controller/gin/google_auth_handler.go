package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	libapi "github.com/kujilabo/cocotola-1.21/lib/api"

	rsliblog "github.com/kujilabo/redstart/lib/log"

	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/domain"
)

type googleAuthParameter struct {
	OrganizationName string `json:"organizationName"`
	SessionState     string `json:"sessionState"`
	ParamState       string `json:"paramState"`
	Code             string `json:"code"`
}

type GoogleUserUsecaseInterface interface {
	GenerateState(context.Context) (string, error)
	Authorize(ctx context.Context, state, code, organizationName string) (*domain.AuthTokenSet, error)
}

type GoogleUserHandler struct {
	googleUserUsecase GoogleUserUsecaseInterface
}

func NewGoogleAuthHandler(googleUserUsecase GoogleUserUsecaseInterface) *GoogleUserHandler {
	return &GoogleUserHandler{
		googleUserUsecase: googleUserUsecase,
	}
}

func (h *GoogleUserHandler) GenerateState(c *gin.Context) {
	ctx := c.Request.Context()
	ctx = rsliblog.WithLoggerName(ctx, loggerKey)
	logger := rsliblog.GetLoggerFromContext(ctx, loggerKey)

	logger.Info("GenerateState")

	state, err := h.googleUserUsecase.GenerateState(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"state": state})
}

func (h *GoogleUserHandler) Authorize(c *gin.Context) {
	ctx := c.Request.Context()
	ctx = rsliblog.WithLoggerName(ctx, loggerKey)
	logger := rsliblog.GetLoggerFromContext(ctx, loggerKey)

	logger.Info("Authorize")

	googleAuthParameter := googleAuthParameter{}
	if err := c.ShouldBindJSON(&googleAuthParameter); err != nil {
		logger.InfoContext(ctx, fmt.Sprintf("invalid parameter. err: %v", err))
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

	// user, err := c.Cookie("auth_user")
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": "auth_user is empty"})
	// 	return
	// }
	if googleAuthParameter.SessionState == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "sessionState is empty"})
		return
	}
	if googleAuthParameter.ParamState == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "paramState is empty"})
		return
	}
	if googleAuthParameter.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "code is empty"})
		return
	}
	if googleAuthParameter.OrganizationName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "organizationName is empty"})
		return
	}
	if googleAuthParameter.SessionState != googleAuthParameter.ParamState {
		c.JSON(http.StatusBadRequest, gin.H{"message": "sessionState and paramState are not equal"})
		return
	}

	authResult, err := h.googleUserUsecase.Authorize(ctx, googleAuthParameter.ParamState, googleAuthParameter.Code, googleAuthParameter.OrganizationName)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthenticated) {
			logger.InfoContext(ctx, fmt.Sprintf("invalid parameter. err: %v", err))
			c.JSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		logger.ErrorContext(ctx, fmt.Sprintf("failed to RegisterStudent. err: %+v", err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, libapi.AuthResponse{
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
		auth.GET("state", googleAuthHandler.GenerateState)
		auth.POST("authorize", googleAuthHandler.Authorize)
		return nil
	}
}
