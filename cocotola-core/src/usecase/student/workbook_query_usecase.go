package student

import (
	"context"
	"errors"

	rsuserdomain "github.com/kujilabo/redstart/user/domain"

	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/service"

	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/domain"
)

type WorkbookQueryService interface {
	RetrieveWorkbookByID(ctx context.Context, workbookID *domain.WorkbookID) (*WorkbookRetrieveModel, error)
}

type WorkbookQueryUsecase struct {
	txManager            service.TransactionManager
	nonTxManager         service.TransactionManager
	workbookQuerySerivce WorkbookQueryService
}

func NewWorkbookQueryUsecase(txManager, nonTxManager service.TransactionManager, workbookQuerySerivce WorkbookQueryService) *WorkbookQueryUsecase {
	return &WorkbookQueryUsecase{
		txManager:            txManager,
		nonTxManager:         nonTxManager,
		workbookQuerySerivce: workbookQuerySerivce,
	}
}

func (u *WorkbookQueryUsecase) FindWorkbooks(ctx context.Context, organizationID *rsuserdomain.OrganizationID, operatorID *rsuserdomain.AppUserID, param *WorkbookFindParameter) (*WorkbookFindResult, error) {
	return &WorkbookFindResult{
		TotalCount: 1,
		Results: []*WorkbookFindWorkbookModel{
			{
				ID:   1,
				Name: "test",
			},
		},
	}, nil

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

func (u *WorkbookQueryUsecase) RetrieveWorkbookByID(ctx context.Context, organizationID *rsuserdomain.OrganizationID, operatorID *rsuserdomain.AppUserID, workbookID int) (*WorkbookRetrieveModel, error) {
	// TODO: check RBAC

	model, err := u.workbookQuerySerivce.RetrieveWorkbookByID(ctx, &domain.WorkbookID{Value: workbookID})
	if err != nil {
		return nil, err
	}

	return model, nil
}

type EnglishSentenceModel struct {
	SrcLang2                  string `json:"srcLang2"`
	SrcAudioContent           string `json:"srcAudioContent"`
	SrcAudioLengthMillisecond int    `json:"SrcAudioLengthMillisecond"`
	SrcText                   string `json:"srcText"`
	DstLang2                  string `json:"dstLang2"`
	DstAudioContent           string `json:"dstAudioContent"`
	DstAudioLengthMillisecond int    `json:"DstAudioLengthMillisecond"`
	DstText                   string `json:"dstText"`
}

type EnglishSentencesModel struct {
	Sentences []*EnglishSentenceModel `json:"sentences"`
}
type EnglishConversationModel struct {
}
type WorkbookRetrieveModel struct {
	ID                  int                       `json:"id"`
	Name                string                    `json:"name"`
	ProblemType         string                    `json:"problmeType"`
	EnglishSentences    *EnglishSentencesModel    `json:"englishSentences,omitempty"`
	EnglishConversation *EnglishConversationModel `json:"englishConversation,omitempty"`
}

type WorkbookFindParameter struct {
	PageNo   int
	PageSize int
}

type WorkbookFindWorkbookModel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type WorkbookFindResult struct {
	TotalCount int                          `json:"totalCount"`
	Results    []*WorkbookFindWorkbookModel `json:"results"`
}

var ErrWorkbookNotFound = errors.New("workbook not found")
