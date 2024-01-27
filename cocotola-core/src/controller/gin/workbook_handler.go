package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/controller/gin/helper"
	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/domain"
	workbookadddomain "github.com/kujilabo/cocotola-1.21/cocotola-core/src/domain/workbookadd"
	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/service"
	studentusecase "github.com/kujilabo/cocotola-1.21/cocotola-core/src/usecase/student"
)

const defaultPageSize = 10

// type WorkbookFindModel struct {
// 	ID   int    `json:"id"`
// 	Name string `json:"name"`
// }

// type WorkbookFindResult struct {
// 	TotalCount int                  `json:"totalCount"`
// 	Results    []*WorkbookFindModel `json:"results"`
// }

//	type Problem struct {
//		Type       string            `json:"type"`
//		Properties map[string]string `json:"properties"`
//	}
//
//	type WorkbookWithProblem struct {
//		ID       int        `json:"id"`
//		Problems []*Problem `json:"problems"`
//	}
type WorkbookQueryUsecase interface {
	FindWorkbooks(ctx context.Context, operator service.OperatorInterface, param *studentusecase.WorkbookFindParameter) (*studentusecase.WorkbookFindResult, error)

	RetrieveWorkbookByID(ctx context.Context, operator service.OperatorInterface, workbookID int) (*studentusecase.WorkbookRetrieveResult, error)
}
type WorkbookCommandUsecase interface {
	AddWorkbook(ctx context.Context, operator service.OperatorInterface, param *workbookadddomain.Parameter) (*domain.WorkbookID, error)
}

type WorkbookHandler struct {
	workbookQueryUsecase   WorkbookQueryUsecase
	workbookCommandUsecase WorkbookCommandUsecase
}

func NewWorkbookHandler(workbookQueryUsecase WorkbookQueryUsecase, workbookCommandUsecase WorkbookCommandUsecase) *WorkbookHandler {
	return &WorkbookHandler{
		workbookQueryUsecase:   workbookQueryUsecase,
		workbookCommandUsecase: workbookCommandUsecase,
	}
}

func (h *WorkbookHandler) FindWorkbooks(c *gin.Context) {
	helper.HandleSecuredFunction(c, func(ctx context.Context, logger *slog.Logger, operator service.OperatorInterface) error {
		param := studentusecase.WorkbookFindParameter{
			PageNo:   1,
			PageSize: defaultPageSize,
		}
		result, err := h.workbookQueryUsecase.FindWorkbooks(ctx, operator, &param)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, result)
		return nil
	}, h.errorHandle)
}

// func (h *WorkbookHandler) toWorkbookFindResultEntity(model *studentusecase.WorkbookFindResult) *WorkbookFindResult {
// 	results := make([]*WorkbookFindModel, len(model.Results))
// 	for i, r := range model.Results {
// 		results[i] = &WorkbookFindModel{ID: r.ID, Name: r.Name}
// 	}

// 	return &WorkbookFindResult{
// 		TotalCount: model.TotalCount,
// 		Results:    results,
// 	}
// }

func (h *WorkbookHandler) RetrieveWorkbookByID(c *gin.Context) {
	helper.HandleSecuredFunction(c, func(ctx context.Context, logger *slog.Logger, operator service.OperatorInterface) error {
		result, err := h.workbookQueryUsecase.RetrieveWorkbookByID(ctx, operator, 1)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, result)
		return nil
	}, h.errorHandle)
}

func (h *WorkbookHandler) AddWorkbook(c *gin.Context) {

}
func (h *WorkbookHandler) UpdateWorkbook(c *gin.Context) {

}

// func (h *WorkbookHandler) toWorkbookRetrieveResultEntity(model *workbookretrievedomain.WorkbookModel) *WorkbookWithProblem {
// 	problems := make([]*Problem, len(model.Problems))
// 	for i, r := range model.Problems {
// 		problems[i] = &Problem{
// 			Type:       r.Type,
// 			Properties: r.Properties,
// 		}
// 	}

// 	return &WorkbookWithProblem{
// 		ID:       model.ID,
// 		Problems: problems,
// 	}
// }

func (h *WorkbookHandler) errorHandle(ctx context.Context, logger *slog.Logger, c *gin.Context, err error) bool {
	// if errors.Is(err, service.ErrAudioNotFound) {
	// 	logger.Warnf("PrivateWorkbookHandler err: %+v", err)
	// 	c.JSON(http.StatusNotFound, gin.H{"message": "Audio not found"})
	// 	return true
	// }
	logger.ErrorContext(ctx, fmt.Sprintf("WorkbookHandler. error: %+v", err))
	return false
}

func NewInitWorkbookRouterFunc(workbookQueryUsecase WorkbookQueryUsecase, workbookCommandUsecase WorkbookCommandUsecase) InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
		workbook := parentRouterGroup.Group("workbook")
		workbookHandler := NewWorkbookHandler(workbookQueryUsecase, workbookCommandUsecase)
		for _, m := range middleware {
			workbook.Use(m)
		}
		workbook.GET("", workbookHandler.FindWorkbooks)
		workbook.GET(":workbookID", workbookHandler.RetrieveWorkbookByID)
		// workbook.POST(":workbookID", privateWorkbookHandler.FindWorkbooks)
		// workbook.GET(":workbookID", privateWorkbookHandler.FindWorkbookByID)
		workbook.PUT(":workbookID", workbookHandler.UpdateWorkbook)
		// workbook.DELETE(":workbookID", privateWorkbookHandler.RemoveWorkbook)
		workbook.POST("", workbookHandler.AddWorkbook)
		return nil
	}
}
