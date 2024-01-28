package student

import (
	"context"

	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/domain"
	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/service"
)

type WorkbookCommandUsecase struct {
	txManager    service.TransactionManager
	nonTxManager service.TransactionManager
}

func NewWorkbookCommandUsecase(txManager, nonTxManager service.TransactionManager) *WorkbookCommandUsecase {
	return &WorkbookCommandUsecase{
		txManager:    txManager,
		nonTxManager: nonTxManager,
	}
}

func (u *WorkbookCommandUsecase) AddWorkbook(ctx context.Context, operator service.OperatorInterface, param *service.WorkbookAddParameter) (*domain.WorkbookID, error) {
	var workbookID *domain.WorkbookID
	if err := u.txManager.Do(ctx, func(rf service.RepositoryFactory) error {
		workbookRepo, err := rf.NewWorkbookRepository(ctx)
		if err != nil {
			return err
		}

		tmpWorkbookID, err := workbookRepo.AddWorkbook(ctx, operator, param)
		if err != nil {
			return err
		}

		workbookID = tmpWorkbookID
		return nil
	}); err != nil {
		return nil, err
	}

	return workbookID, nil
}

func (u *WorkbookCommandUsecase) UpdateWorkbook(ctx context.Context, operator service.OperatorInterface, workbookID *domain.WorkbookID, version int, param *service.WorkbookUpdateParameter) error {
	if err := u.txManager.Do(ctx, func(rf service.RepositoryFactory) error {
		workbookRepo, err := rf.NewWorkbookRepository(ctx)
		if err != nil {
			return err
		}

		if err := workbookRepo.UpdateWorkbook(ctx, operator, workbookID, version, param); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}
