package student

import (
	"context"

	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/domain"
	workbookadddomain "github.com/kujilabo/cocotola-1.21/cocotola-core/src/domain/workbookadd"
	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/service"
)

type WorkbookCommandUsecase struct {
}

func NewWorkbookCommandUsecase() *WorkbookCommandUsecase {
	return &WorkbookCommandUsecase{}
}

func (u *WorkbookCommandUsecase) AddWorkbook(ctx context.Context, operator service.OperatorInterface, param *workbookadddomain.Parameter) (*domain.WorkbookID, error) {
	return nil, nil
}
