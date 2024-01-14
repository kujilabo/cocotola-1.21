//go:build small

package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	libapi "github.com/kujilabo/cocotola-1.21/lib/api"
	libconfig "github.com/kujilabo/cocotola-1.21/lib/config"

	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/config"
	handler "github.com/kujilabo/cocotola-1.21/cocotola-core/src/controller/gin"
	handlermock "github.com/kujilabo/cocotola-1.21/cocotola-core/src/controller/gin/mocks"
	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/domain/workbookfind"
	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/service"
	servicemock "github.com/kujilabo/cocotola-1.21/cocotola-core/src/service/mocks"
)

var (
	anythingOfContext = mock.MatchedBy(func(_ context.Context) bool { return true })
	corsConfig        cors.Config
	appConfig         *config.AppConfig
	debugConfig       *libconfig.DebugConfig
	// authTokenManager  auth.AuthTokenManager
)

func init() {
	corsConfig = cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"*"},
		AllowHeaders:    []string{"*"},
	}
	appConfig = &config.AppConfig{
		Name:        "test",
		HTTPPort:    8080,
		MetricsPort: 8081,
	}
	debugConfig = &libconfig.DebugConfig{
		Gin:  false,
		Wait: false,
	}
}

func initWorkbookRouter(t *testing.T, ctx context.Context, cocotolaAuthClient service.CocotolaAuthClient, workbokUsecase handler.WorkbookUsecaseInterface) *gin.Engine {
	t.Helper()
	fn := handler.NewInitWorkbookRouterFunc(workbokUsecase)

	initPublicRouterFunc := []handler.InitRouterGroupFunc{}
	initPrivateRouterFunc := []handler.InitRouterGroupFunc{fn}

	router := gin.New()
	err := handler.InitRouter(ctx, router, cocotolaAuthClient, initPublicRouterFunc, initPrivateRouterFunc, corsConfig, debugConfig, appConfig.Name)
	require.NoError(t, err)

	return router
}

func TestWorkbookHandler_FindWorkbook(t *testing.T) {
	ctx := context.Background()

	type args struct {
		authorizationHeader string
	}
	type outputs struct {
		statusCode  int
		totalCount  int
		resultsLen  int
		resultID0   int
		resultName0 string
	}
	tests := []struct {
		name    string
		args    args
		outputs outputs
	}{
		{
			name: "authorization header is not specified",
			outputs: outputs{
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name: "authorization header is invalid",
			args: args{
				authorizationHeader: "Bearer INVALID_TOKEN",
			},
			outputs: outputs{
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name: "authorization header is valid",
			args: args{
				authorizationHeader: "Bearer VALID_TOKEN",
			},
			outputs: outputs{
				statusCode:  http.StatusOK,
				totalCount:  789,
				resultsLen:  1,
				resultID0:   135,
				resultName0: "WORKBOOK_NAME",
			},
		},
	}

	appUserID := appUserID(t, 123)
	organizaionID := organizationID(t, 456)
	apiUserInfo := libapi.AppUserInfoResponse{
		AppUserID:      123,
		OrganizationID: 456,
		LoginID:        "LOGIN_ID",
		Username:       "USERNAME",
	}
	workbookFindResult := workbookfind.Result{
		TotalCount: 789,
		Results: []*workbookfind.WorkbookModel{
			{
				ID:   135,
				Name: "WORKBOOK_NAME",
			},
		},
	}
	cocotolaAuthClient := new(servicemock.CocotolaAuthClient)
	cocotolaAuthClient.On("RetrieveUserInfo", anythingOfContext, "VALID_TOKEN").Return(&apiUserInfo, nil)
	cocotolaAuthClient.On("RetrieveUserInfo", anythingOfContext, "INVALID_TOKEN").Return(nil, errors.New("invalid token"))
	workbookUsecase := new(handlermock.WorkbookUsecaseInterface)
	workbookUsecase.On("FindWorkbooks", anythingOfContext, organizaionID, appUserID, mock.Anything).Return(&workbookFindResult, nil)
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// given
			r := initWorkbookRouter(t, ctx, cocotolaAuthClient, workbookUsecase)
			w := httptest.NewRecorder()

			// when
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/v1/workbook", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", tt.args.authorizationHeader)
			r.ServeHTTP(w, req)
			respBytes := readBytes(t, w.Body)

			// - check the status code
			assert.Equal(t, tt.outputs.statusCode, w.Code)

			if w.Code != http.StatusOK {
				assert.Len(t, respBytes, 0)
				return
			}

			jsonObj := parseJSON(t, respBytes)

			// - check the totalCount
			totalCountExpr := parseExpr(t, "$.totalCount")
			totalCount := totalCountExpr.Get(jsonObj)
			assert.Equal(t, int64(tt.outputs.totalCount), totalCount[0])

			// - check the results
			resultsExpr := parseExpr(t, "$.results")
			results := resultsExpr.Get(jsonObj)
			assert.Equal(t, tt.outputs.resultsLen, len(results))

			// - check the results[0].id
			resultID0Expr := parseExpr(t, "$.results[0].id")
			resultID0 := resultID0Expr.Get(jsonObj)
			assert.Equal(t, int64(tt.outputs.resultID0), resultID0[0])

			// - check the results[0].name
			resultName0Expr := parseExpr(t, "$.results[0].name")
			resultName0 := resultName0Expr.Get(jsonObj)
			assert.Equal(t, tt.outputs.resultName0, resultName0[0])
		})
	}
}
