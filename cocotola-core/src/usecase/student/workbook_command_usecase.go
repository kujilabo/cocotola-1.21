package student

import (
	"context"

	rsuserdomain "github.com/kujilabo/redstart/user/domain"

	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/domain"
	workbookadddomain "github.com/kujilabo/cocotola-1.21/cocotola-core/src/domain/workbookadd"
)

type WorkbookCommandUsecase struct {
}

func NewWorkbookCommandUsecase() *WorkbookCommandUsecase {
	return &WorkbookCommandUsecase{}
}

func (u *WorkbookCommandUsecase) AddWorkbook(ctx context.Context, organizationID *rsuserdomain.OrganizationID, operatorID *rsuserdomain.AppUserID, param *workbookadddomain.Parameter) (*domain.WorkbookID, error) {
	return nil, nil
}
