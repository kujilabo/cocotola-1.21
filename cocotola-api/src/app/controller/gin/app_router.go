package handler

import (
	"context"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/kujilabo/cocotola-1.21/cocotola-api/src/app/config"
	"github.com/kujilabo/cocotola-1.21/cocotola-api/src/app/controller/gin/middleware"
	studentusecase "github.com/kujilabo/cocotola-1.21/cocotola-api/src/app/usecase/student"
	liblog "github.com/kujilabo/cocotola-1.21/lib/log"
	rsliblog "github.com/kujilabo/redstart/lib/log"
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
func NewInitWorkbookRouterFunc(studentUsecaseWorkbook studentusecase.StudentUsecaseWorkbook) InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
		// workbook := parentRouterGroup.Group("private/workbook")
		// privateWorkbookHandler := NewPrivateWorkbookHandler(studentUsecaseWorkbook)
		// for _, m := range middleware {
		// 	workbook.Use(m)
		// }
		// workbook.POST(":workbookID", privateWorkbookHandler.FindWorkbooks)
		// workbook.GET(":workbookID", privateWorkbookHandler.FindWorkbookByID)
		// workbook.PUT(":workbookID", privateWorkbookHandler.UpdateWorkbook)
		// workbook.DELETE(":workbookID", privateWorkbookHandler.RemoveWorkbook)
		// workbook.POST("", privateWorkbookHandler.AddWorkbook)
		return nil
	}
}

func NewAppRouter(
	ctx context.Context,
	initPublicRouterFunc []InitRouterGroupFunc,
	// initPrivateRouterFunc []InitRouterGroupFunc, initPluginRouterFunc []InitRouterGroupFunc,
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

	if debugConfig.Wait {
		router.Use(middleware.NewWaitMiddleware())
	}

	v1 := router.Group("v1")
	{
		v1.Use(otelgin.Middleware(appConfig.Name))
		v1.Use(middleware.NewTraceLogMiddleware(appConfig.Name))

		for _, fn := range initPublicRouterFunc {
			if err := fn(v1); err != nil {
				return nil, err
			}
		}
	}

	return router, nil
}
