package handler

import (
	"context"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	rsliblog "github.com/kujilabo/redstart/lib/log"

	libmiddleware "github.com/kujilabo/cocotola-1.21/lib/controller/gin/middleware"
	liblog "github.com/kujilabo/cocotola-1.21/lib/log"

	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/app/config"
	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/app/controller/gin/middleware"
	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/app/service"
)

// type NewIteratorFunc func(ctx context.Context, workbookID appD.WorkbookID, problemType appD.ProblemTypeName, reader io.Reader) (appS.ProblemAddParameterIterator, error)

type InitRouterGroupFunc func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error

func NewInitTestRouterFunc() InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
		test := parentRouterGroup.Group("test")
		for _, m := range middleware {
			test.Use(m)
		}
		test.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})
		return nil
	}
}

func NewAppRouter(
	ctx context.Context, cocotolaAuthClient service.CocotolaAuthClient,
	initPublicRouterFunc []InitRouterGroupFunc,
	initPrivateRouterFunc []InitRouterGroupFunc,
	// initPluginRouterFunc []InitRouterGroupFunc,
	//authTokenManager service.AuthTokenManager,
	corsConfig cors.Config, appConfig *config.AppConfig,
	// authConfig *config.AuthConfig,
	debugConfig *config.DebugConfig) (*gin.Engine, error) {
	logger := rsliblog.GetLoggerFromContext(ctx, liblog.AppControllerLoggerContextKey)

	router := gin.New()
	router.Use(cors.New(corsConfig))
	router.Use(gin.Recovery())

	// if debugConfig.GinMode {
	// 	// router.Use(ginlog.Middleware(ginlog.DefaultConfig))
	// }
	router.Use(sloggin.New(logger))
	authMiddleware := middleware.NewAuthMiddleware(cocotolaAuthClient)

	if debugConfig.Wait {
		router.Use(libmiddleware.NewWaitMiddleware())
	}

	v1 := router.Group("v1")
	{
		v1.Use(otelgin.Middleware(appConfig.Name))
		v1.Use(libmiddleware.NewTraceLogMiddleware(appConfig.Name))

		for _, fn := range initPublicRouterFunc {
			if err := fn(v1); err != nil {
				return nil, err
			}
		}
		for _, fn := range initPrivateRouterFunc {
			if err := fn(v1, authMiddleware); err != nil {
				return nil, err
			}
		}
	}

	return router, nil
}
