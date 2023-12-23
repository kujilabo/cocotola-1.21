package handler

// import (
// 	"github.com/gin-gonic/gin"

// 	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/service"
// 	liblog "github.com/kujilabo/cocotola-1.21/lib/log"
// 	rsliblog "github.com/kujilabo/redstart/lib/log"
// )

// type AuthHandler interface {
// 	RefreshToken(c *gin.Context)
// }

// type authHandler struct {
// 	authTokenManager service.AuthTokenManager
// }

// func NewAuthHandler(authTokenManager service.AuthTokenManager) AuthHandler {
// 	return &authHandler{
// 		authTokenManager: authTokenManager,
// 	}
// }

// func (h *authHandler) RefreshToken(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	logger := rsliblog.GetLoggerFromContext(ctx, liblog.AppControllerLoggerContextKey)
// 	logger.Info("Authorize")
// 	refreshTokenParameter := RefreshTokenParameter{}
// 	if err := c.BindJSON(&refreshTokenParameter); err != nil {
// 		return
// 	}

// 	// token, err := h.authTokenManager.RefreshToken(ctx, refreshTokenParameter.RefreshToken)
// 	// if err != nil {
// 	// 	c.Status(http.StatusBadRequest)
// 	// 	return
// 	// }

// 	// logger.Info("Authorize OK")
// 	// c.JSON(http.StatusOK, AuthResponse{
// 	// 	AccessToken: token,
// 	// })
// }
