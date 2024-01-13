package initialize

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	libconfig "github.com/kujilabo/cocotola-1.21/lib/config"
	rslibconfig "github.com/kujilabo/redstart/lib/config"
	rsuserservice "github.com/kujilabo/redstart/user/service"
	"gorm.io/gorm"

	controller "github.com/kujilabo/cocotola-1.21/cocotola-core/src/controller/gin"
	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/gateway"
	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/service"
	studentusecase "github.com/kujilabo/cocotola-1.21/cocotola-core/src/usecase/student"
)

// const readHeaderTimeout = time.Duration(30) * time.Second
const authClientTimeout = time.Duration(5) * time.Second

func InitTransactionManager(db *gorm.DB, rff gateway.RepositoryFactoryFunc) service.TransactionManager {
	appTransactionManager, err := gateway.NewTransactionManager(db, rff)
	if err != nil {
		panic(err)
	}

	return appTransactionManager
}

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

func InitAppServer(ctx context.Context, parentRouterGroup gin.IRouter, corsConfig *rslibconfig.CORSConfig, debugConfig *libconfig.DebugConfig, appName string, transactionManager service.TransactionManager, rsrf rsuserservice.RepositoryFactory) error {
	// cors
	gincorsConfig := rslibconfig.InitCORS(corsConfig)
	// httpClient := http.Client{
	// 	Timeout:   time.Duration(authConfig.APITimeoutSec) * time.Second,
	// 	Transport: otelhttp.NewTransport(http.DefaultTransport),
	// }
	studentUsecaseWorkbook := studentusecase.NewStudentUsecaseWorkbook(transactionManager)
	privateRouterGroupFunc := []controller.InitRouterGroupFunc{
		controller.NewInitWorkbookRouterFunc(studentUsecaseWorkbook),
		// controller.NewInitProblemRouterFunc(studentUsecaseProblem, newIteratorFunc),
		// controller.NewInitStudyRouterFunc(studentUseCaseStudy),
		// controller.NewInitAudioRouterFunc(studentUsecaseAudio),
		// controller.NewInitStatRouterFunc(studentUsecaseStat),
	}

	publicRouterGroupFunc := []controller.InitRouterGroupFunc{
		controller.NewInitTestRouterFunc(),
	}
	cocotolaAuthClient := gateway.NewCocotolaAuthClient(authClientTimeout)

	if err := controller.InitRouter(ctx, parentRouterGroup, cocotolaAuthClient, publicRouterGroupFunc, privateRouterGroupFunc, gincorsConfig, debugConfig, appName); err != nil {
		return err
	}

	return nil
}
