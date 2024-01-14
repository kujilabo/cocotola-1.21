package student

import (
	"context"
	"errors"

	rsuserdomain "github.com/kujilabo/redstart/user/domain"

	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/service"

	workbookfinddomain "github.com/kujilabo/cocotola-1.21/cocotola-core/src/domain/workbookfind"
)

type StudentUsecaseWorkbookInterface interface {
	// FindWorkbooks(ctx context.Context, organizationID *rsuserdomain.OrganizationID, operatorID *rsuserdomain.AppUserID, param *workbookfinddomain.Parameter) (*workbookfinddomain.Result, error)
}

type StudentUsecaseWorkbook struct {
	txManager    service.TransactionManager
	nonTxManager service.TransactionManager
}

func NewStudentUsecaseWorkbook(txManager service.TransactionManager, nonTxManager service.TransactionManager) *StudentUsecaseWorkbook {
	return &StudentUsecaseWorkbook{
		txManager:    txManager,
		nonTxManager: nonTxManager,
	}
}

func (u *StudentUsecaseWorkbook) FindWorkbooks(ctx context.Context, organizationID *rsuserdomain.OrganizationID, operatorID *rsuserdomain.AppUserID, param *workbookfinddomain.Parameter) (*workbookfinddomain.Result, error) {
	return nil, errors.New("")
	// var result domain.WorkbookSearchResult
	// fn := func(student service.Student) error {
	// 	condition, err := domain.NewWorkbookSearchCondition(DefaultPageNo, DefaultPageSize, []userD.SpaceID{})
	// 	if err != nil {
	// 		return rserrors.Errorf("service.NewWorkbookSearchCondition. err: %w", err)
	// 	}

	// 	tmpResult, err := student.FindWorkbooksFromPersonalSpace(ctx, condition)
	// 	if err != nil {
	// 		return rserrors.Errorf("student.FindWorkbooksFromPersonalSpace. err: %w", err)
	// 	}

	// 	result = tmpResult
	// 	return nil
	// }

	// if err := u.studentHandle(ctx, organizationID, operatorID, fn); err != nil {
	// 	return nil, err
	// }

	// return result, nil
}
