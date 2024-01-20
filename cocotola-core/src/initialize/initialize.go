package initialize

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	libconfig "github.com/kujilabo/cocotola-1.21/lib/config"
	rslibconfig "github.com/kujilabo/redstart/lib/config"
	rsuserservice "github.com/kujilabo/redstart/user/service"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/config"
	controller "github.com/kujilabo/cocotola-1.21/cocotola-core/src/controller/gin"
	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/controller/gin/middleware"
	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/gateway"
	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/service"
	studentusecase "github.com/kujilabo/cocotola-1.21/cocotola-core/src/usecase/student"
)

// const readHeaderTimeout = time.Duration(30) * time.Second
const authClientTimeout = time.Duration(5) * time.Second

// func InitTransactionManager(db *gorm.DB, rff gateway.RepositoryFactoryFunc) service.TransactionManager {
// 	appTransactionManager, err := gateway.NewTransactionManager(db, rff)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return appTransactionManager
// }

// type systemOwnerByOrganizationName struct {
// }

// func (s systemOwnerByOrganizationName) Get(ctx context.Context, rf service.RepositoryFactory, organizationName string) (*rsuserservice.SystemOwner, error) {
// 	rsrf, err := rf.NewRedstartRepositoryFactory(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	systemAdmin, err := rsuserservice.NewSystemAdmin(ctx, rsrf)
// 	if err != nil {
// 		return nil, err
// 	}

// 	systemOwner, err := systemAdmin.FindSystemOwnerByOrganizationName(ctx, organizationName)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return systemOwner, nil
// }

func InitAppServer(ctx context.Context, parentRouterGroup gin.IRouter, authAPIConfig config.AuthAPIonfig, corsConfig *rslibconfig.CORSConfig, debugConfig *libconfig.DebugConfig, appName string, txManager service.TransactionManager, nonTxManager service.TransactionManager, rsrf rsuserservice.RepositoryFactory) error {
	// cors
	gincorsConfig := rslibconfig.InitCORS(corsConfig)
	httpClient := http.Client{
		Timeout:   authClientTimeout,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	studentUsecaseWorkbook := studentusecase.NewStudentUsecaseWorkbook(txManager, nonTxManager)
	privateRouterGroupFunc := []controller.InitRouterGroupFunc{
		controller.NewInitPrivateWorkbookRouterFunc(studentUsecaseWorkbook),
		controller.NewInitWorkbookRouterFunc(studentUsecaseWorkbook),
		// controller.NewInitProblemRouterFunc(studentUsecaseProblem, newIteratorFunc),
		// controller.NewInitStudyRouterFunc(studentUseCaseStudy),
		// controller.NewInitAudioRouterFunc(studentUsecaseAudio),
		// controller.NewInitStatRouterFunc(studentUsecaseStat),
	}
	authEndpoint, err := url.Parse(authAPIConfig.Endpoint)
	if err != nil {
		return err
	}

	publicRouterGroupFunc := []controller.InitRouterGroupFunc{
		controller.NewInitTestRouterFunc(),
	}
	cocotolaAuthClient := gateway.NewCocotolaAuthClient(&httpClient, authEndpoint, authAPIConfig.Username, authAPIConfig.Password)
	authMiddleware := middleware.NewAuthMiddleware(cocotolaAuthClient)

	if err := controller.InitRouter(ctx, parentRouterGroup, authMiddleware, publicRouterGroupFunc, privateRouterGroupFunc, gincorsConfig, debugConfig, appName); err != nil {
		return err
	}

	return nil
}
