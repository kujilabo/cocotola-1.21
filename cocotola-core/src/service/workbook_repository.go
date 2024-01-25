package service

import (
	"context"
	"errors"

	rsuserdomain "github.com/kujilabo/redstart/user/domain"

	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/domain"
)

var ErrWorkbookAlreadyExists = errors.New("Workbook already exists")

type WorkbookAddParameter struct {
	Name        string
	ProblemType string
	Lang2       string
	Description string
	Content     string
}

type WorkbookUpdateParameter struct {
	Name        string
	Description string
	Content     string
}

type AppUserIDInterface interface {
	AppUserID() *rsuserdomain.AppUserID
	OrganizationID() *rsuserdomain.OrganizationID
	// LoginID() string
	// Username() string
}

type WorkbookRepository interface {
	AddWorkbook(ctx context.Context, operator AppUserIDInterface, param WorkbookAddParameter) (*domain.WorkbookID, error)

	UpdateWorkbook(ctx context.Context, operator AppUserIDInterface, workbookID *domain.WorkbookID, version int, param WorkbookUpdateParameter) error
}
