package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/usecase"
	liblog "github.com/kujilabo/cocotola-1.21/lib/log"
	rsliblog "github.com/kujilabo/redstart/lib/log"
)

type GoogleUserHandler interface {
	Authorize(c *gin.Context)
}

type googleUserHandler struct {
	googleUserUsecase usecase.GoogleUserUsecaseInterface
}

func NewGoogleAuthHandler(googleUserUsecase usecase.GoogleUserUsecaseInterface) GoogleUserHandler {
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
