package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	rsuserdomain "github.com/kujilabo/redstart/user/domain"

	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/controller/gin/helper"
	workbookfinddomain "github.com/kujilabo/cocotola-1.21/cocotola-core/src/domain/workbookfind"
)

const defaultPageSize = 10

type WorkbookFindModel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type WorkbookFindResult struct {
	TotalCount int                  `json:"totalCount"`
	Results    []*WorkbookFindModel `json:"results"`
}

type WorkbookUsecaseInterface interface {
	FindWorkbooks(ctx context.Context, organizationID *rsuserdomain.OrganizationID, operatorID *rsuserdomain.AppUserID, param *workbookfinddomain.Parameter) (*workbookfinddomain.Result, error)
}

type WorkbookHandler struct {
	workbookusecase WorkbookUsecaseInterface
}

func NewWorkbookHandler(workbookusecase WorkbookUsecaseInterface) *WorkbookHandler {
	return &WorkbookHandler{
		workbookusecase: workbookusecase,
	}
}

func (h *WorkbookHandler) FindWorkbooks(c *gin.Context) {
	helper.HandleSecuredFunction(c, func(ctx context.Context, logger *slog.Logger, organizationID *rsuserdomain.OrganizationID, operatorID *rsuserdomain.AppUserID) error {
		param := workbookfinddomain.Parameter{
			PageNo:   1,
			PageSize: defaultPageSize,
		}
		result, err := h.workbookusecase.FindWorkbooks(ctx, organizationID, operatorID, &param)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, h.toWorkbookFindResultEntity(result))
		return nil
	}, h.errorHandle)
}

func (h *WorkbookHandler) toWorkbookFindResultEntity(model *workbookfinddomain.Result) *WorkbookFindResult {
	results := make([]*WorkbookFindModel, len(model.Results))
	for i, r := range model.Results {
		results[i] = &WorkbookFindModel{ID: r.ID, Name: r.Name}
	}

	return &WorkbookFindResult{
		TotalCount: model.TotalCount,
		Results:    results,
	}
}

func (h *WorkbookHandler) RetrieveWorkbookByID(c *gin.Context) {

}

func (h *WorkbookHandler) errorHandle(ctx context.Context, logger *slog.Logger, c *gin.Context, err error) bool {
	// if errors.Is(err, service.ErrAudioNotFound) {
	// 	logger.Warnf("PrivateWorkbookHandler err: %+v", err)
	// 	c.JSON(http.StatusNotFound, gin.H{"message": "Audio not found"})
	// 	return true
	// }
	logger.ErrorContext(ctx, fmt.Sprintf("error:%v", err))
	return false
}

func NewInitWorkbookRouterFunc(workbookUsecase WorkbookUsecaseInterface) InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
		workbook := parentRouterGroup.Group("workbook")
		workbookHandler := NewWorkbookHandler(workbookUsecase)
		for _, m := range middleware {
			workbook.Use(m)
		}
		workbook.GET("", workbookHandler.FindWorkbooks)
		// workbook.POST(":workbookID", privateWorkbookHandler.FindWorkbooks)
		// workbook.GET(":workbookID", privateWorkbookHandler.FindWorkbookByID)
		// workbook.PUT(":workbookID", privateWorkbookHandler.UpdateWorkbook)
		// workbook.DELETE(":workbookID", privateWorkbookHandler.RemoveWorkbook)
		// workbook.POST("", privateWorkbookHandler.AddWorkbook)
		return nil
	}
}
