package handler

// import (
// 	"context"
// 	"net/http"

// 	"github.com/gin-contrib/cors"
// 	"github.com/gin-gonic/gin"
// 	sloggin "github.com/samber/slog-gin"
// 	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

// 	rsliblog "github.com/kujilabo/redstart/lib/log"

// 	liblog "github.com/kujilabo/cocotola-1.21/lib/log"

// 	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/config"
// 	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/controller/gin/middleware"
// 	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/service"
// 	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/usecase"
// )

// // type NewIteratorFunc func(ctx context.Context, workbookID appD.WorkbookID, problemType appD.ProblemTypeName, reader io.Reader) (appS.ProblemAddParameterIterator, error)

// type InitRouterGroupFunc func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error

// func NewInitTestRouterFunc() InitRouterGroupFunc {
// 	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
// 		test := parentRouterGroup.Group("test")
// 		for _, m := range middleware {
// 			test.Use(m)
// 		}
// 		test.GET("/ping", func(c *gin.Context) {
// 			c.String(http.StatusOK, "pong")
// 		})
// 		return nil
// 	}
// }

// func NewInitAuthRouterFunc(googleUserUsecase usecase.GoogleUserUsecase, authTokenManager service.AuthTokenManager) InitRouterGroupFunc {
// 	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
// 		auth := parentRouterGroup.Group("auth")
// 		for _, m := range middleware {
// 			auth.Use(m)
// 		}
// 		// googleUserUsecase := authU.NewGoogleUserUsecase(db, googleAuthClient, authTokenManager, registerAppUserCallback)
// 		// guestUserUsecase := authU.NewGuestUserUsecase(authTokenManager)
// 		authHandler := NewAuthHandler(authTokenManager)
// 		googleAuthHandler := NewGoogleAuthHandler(googleUserUsecase)
// 		// guestAuthHandler := authH.NewGuestAuthHandler(guestUserUsecase)
// 		auth.POST("google/authorize", googleAuthHandler.Authorize)
// 		auth.POST("refresh_token", authHandler.RefreshToken)
// 		return nil
// 	}
// }
// func NewAppRouter(
// 	ctx context.Context,
// 	initPublicRouterFunc []InitRouterGroupFunc,
// 	initPrivateRouterFunc []InitRouterGroupFunc,
// 	// initPluginRouterFunc []InitRouterGroupFunc,
// 	//authTokenManager service.AuthTokenManager,
// 	corsConfig cors.Config, appConfig *config.AppConfig,
// 	// authConfig *config.AuthConfig,
// 	debugConfig *config.DebugConfig) (*gin.Engine, error) {
// 	logger := rsliblog.GetLoggerFromContext(ctx, liblog.AppControllerLoggerContextKey)

// 	router := gin.New()
// 	router.Use(cors.New(corsConfig))
// 	// router.Use(gin.Recovery())
// 	router.Use(sloggin.New(logger))

// 	if debugConfig.Wait {
// 		router.Use(middleware.NewWaitMiddleware())
// 	}

// 	v1 := router.Group("v1")
// 	{
// 		v1.Use(otelgin.Middleware(appConfig.Name))
// 		v1.Use(middleware.NewTraceLogMiddleware(appConfig.Name))

// 		for _, fn := range initPublicRouterFunc {
// 			if err := fn(v1); err != nil {
// 				return nil, err
// 			}
// 		}
// 		// for _, fn := range initPrivateRouterFunc {
// 		// 	if err := fn(v1, authMiddleware); err != nil {
// 		// 		return nil, err
// 		// 	}
// 		// }
// 	}

// 	return router, nil
// }
