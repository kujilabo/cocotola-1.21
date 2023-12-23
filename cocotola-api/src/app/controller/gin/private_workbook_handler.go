package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	studentusecase "github.com/kujilabo/cocotola-1.21/cocotola-api/src/app/usecase/student"
	rsliblog "github.com/kujilabo/redstart/lib/log"
)

// type PrivateWorkbookHandler interface {
// 	FindWorkbooks(c *gin.Context)
// 	// FindWorkbookByID(c *gin.Context)
// 	// AddWorkbook(c *gin.Context)
// 	// UpdateWorkbook(c *gin.Context)
// 	// RemoveWorkbook(c *gin.Context)
// }

type PrivateWorkbookHandler struct {
	// repository             gateway.Repository
	studentUsecaseWorkbook studentusecase.StudentUsecaseWorkbook
}

func NewPrivateWorkbookHandler(studentUsecaseWorkbook studentusecase.StudentUsecaseWorkbook) *PrivateWorkbookHandler {
	return &PrivateWorkbookHandler{
		studentUsecaseWorkbook: studentUsecaseWorkbook,
	}
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
