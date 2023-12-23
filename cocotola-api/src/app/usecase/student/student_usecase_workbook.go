package student

import (
	"context"
	"errors"

	"github.com/kujilabo/cocotola-1.21/cocotola-api/src/app/service"
	rsuserdomain "github.com/kujilabo/redstart/user/domain"
)

type StudentUsecaseWorkbook struct {
	transactionManager service.TransactionManager
}

func NewStudentUsecaseWorkbook(transactionManager service.TransactionManager) *StudentUsecaseWorkbook {
	return &StudentUsecaseWorkbook{
		transactionManager: transactionManager,
	}
}
func (u *StudentUsecaseWorkbook) FindWorkbooks(ctx context.Context, organizationID rsuserdomain.OrganizationID, operatorID rsuserdomain.AppUserID) error {
	return errors.New("")
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
