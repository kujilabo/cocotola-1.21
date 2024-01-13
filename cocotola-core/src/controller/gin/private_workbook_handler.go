package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	rsliblog "github.com/kujilabo/redstart/lib/log"
	rsuserdomain "github.com/kujilabo/redstart/user/domain"

	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/controller/gin/helper"
	studentusecase "github.com/kujilabo/cocotola-1.21/cocotola-core/src/usecase/student"
)

type PrivateWorkbookHandler struct {
	// repository             gateway.Repository
	studentUsecaseWorkbook studentusecase.StudentUsecaseWorkbookInterface
}

func NewPrivateWorkbookHandler(studentUsecaseWorkbook studentusecase.StudentUsecaseWorkbookInterface) *PrivateWorkbookHandler {
	return &PrivateWorkbookHandler{
		studentUsecaseWorkbook: studentUsecaseWorkbook,
	}
}

func (h *PrivateWorkbookHandler) Test(c *gin.Context) {
	helper.HandleSecuredFunction(c, func(ctx context.Context, logger *slog.Logger, organizationID *rsuserdomain.OrganizationID, operatorID *rsuserdomain.AppUserID) error {
		logger.InfoContext(ctx, "TEST")
		c.Status(http.StatusOK)
		return nil
	}, h.errorHandle)
}

// FindWorkbooks godoc
// @Summary Find workbooks
// @Produce json
// @Success 200 {object} entity.WorkbookSearchResponse
// @Failure 400
// @Router /v1/private/workbook/search [post]
func (h *PrivateWorkbookHandler) FindWorkbooks(c *gin.Context) {
	ctx := c.Request.Context()
	logger := rsliblog.GetLoggerFromContext(ctx, HandlerContextKey)

	logger.Info("FindWorkbooks")

	id := c.Param("workbookID")
	if id != "search" {
		c.Status(http.StatusNotFound)
		return
	}

	// controllerhelper.HandleSecuredFunction(c, func(organizationID userD.OrganizationID, operatorID userD.AppUserID) error {
	// 	result, err := h.studentUsecaseWorkbook.FindWorkbooks(ctx, organizationID, operatorID)
	// 	if err != nil {
	// 		return liberrors.Errorf("h.studentUsecaseWorkbook.FindWorkbooks. err: %w", err)
	// 	}

	// 	response, err := converter.ToWorkbookSearchResponse(result)
	// 	if err != nil {
	// 		return liberrors.Errorf("converter.ToWorkbookSearchResponse. err: %w", err)
	// 	}
	// 	c.JSON(http.StatusOK, response)
	// 	return nil
	// }, h.errorHandle)
}

func (h *PrivateWorkbookHandler) errorHandle(ctx context.Context, logger *slog.Logger, c *gin.Context, err error) bool {
	// if errors.Is(err, service.ErrAudioNotFound) {
	// 	logger.Warnf("PrivateWorkbookHandler err: %+v", err)
	// 	c.JSON(http.StatusNotFound, gin.H{"message": "Audio not found"})
	// 	return true
	// }
	logger.ErrorContext(ctx, fmt.Sprintf("error:%v", err))
	return false
}

func NewInitWorkbookRouterFunc(studentUsecaseWorkbook studentusecase.StudentUsecaseWorkbookInterface) InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
		workbook := parentRouterGroup.Group("private/workbook")
		privateWorkbookHandler := NewPrivateWorkbookHandler(studentUsecaseWorkbook)
		for _, m := range middleware {
			workbook.Use(m)
		}
		workbook.GET("test", privateWorkbookHandler.Test)
		// workbook.POST(":workbookID", privateWorkbookHandler.FindWorkbooks)
		// workbook.GET(":workbookID", privateWorkbookHandler.FindWorkbookByID)
		// workbook.PUT(":workbookID", privateWorkbookHandler.UpdateWorkbook)
		// workbook.DELETE(":workbookID", privateWorkbookHandler.RemoveWorkbook)
		// workbook.POST("", privateWorkbookHandler.AddWorkbook)
		return nil
	}
}
