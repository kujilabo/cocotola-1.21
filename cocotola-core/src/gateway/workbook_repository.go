package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"

	rsliberrors "github.com/kujilabo/redstart/lib/errors"
	rslibgateway "github.com/kujilabo/redstart/lib/gateway"

	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/domain"
	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/service"
	libapi "github.com/kujilabo/cocotola-1.21/lib/api"
)

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

func FromEnglishSentenceModel(model *EnglishSentencesModel) ([]byte, error) {
	return json.Marshal(model)
}

func ToEnglishSentenceModel(content []byte) (*EnglishSentencesModel, error) {
	model := EnglishSentencesModel{}
	if err := json.Unmarshal(content, &model); err != nil {
		return nil, err
	}

	return &model, nil
}

type WorkbookEntity struct {
	ID             int
	Version        int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CreatedBy      int
	UpdatedBy      int
	OrganizationID int
	Name           string
	ProblemType    string
	Lang2          string
	Description    string
	Content        string
}

func (e *WorkbookEntity) TableName() string {
	return "workbook"
}

func (e *WorkbookEntity) ToModel() (*libapi.WorkbookRetrieveResult, error) {
	if e.ProblemType == "english_sentences" {
		gwEngSentences, err := ToEnglishSentenceModel([]byte(e.Content))
		if err != nil {
			return nil, err
		}

		sentences := make([]*libapi.EnglishSentenceModel, len(gwEngSentences.Sentences))
		for i := range gwEngSentences.Sentences {
			gwEngSentence := gwEngSentences.Sentences[i]
			sentences[i] = &libapi.EnglishSentenceModel{
				SrcLang2:                  gwEngSentence.SrcLang2,
				SrcAudioContent:           gwEngSentence.SrcAudioContent,
				SrcAudioLengthMillisecond: gwEngSentence.SrcAudioLengthMillisecond,
				SrcText:                   gwEngSentence.SrcText,
				DstLang2:                  gwEngSentence.DstLang2,
				DstAudioContent:           gwEngSentence.DstAudioContent,
				DstAudioLengthMillisecond: gwEngSentence.DstAudioLengthMillisecond,
				DstText:                   gwEngSentence.DstText,
			}
		}

		return &libapi.WorkbookRetrieveResult{
			ID:          e.ID,
			Version:     e.Version,
			Name:        e.Name,
			ProblemType: e.ProblemType,
			EnglishSentences: &libapi.EnglishSentencesModel{
				Sentences: sentences,
			},
		}, nil
	}

	return nil, errors.New("NOT SUPPORTED")
}

type workbookRepository struct {
	db *gorm.DB
}

func NewWorkbookRepository(db *gorm.DB) service.WorkbookRepository {
	return &workbookRepository{
		db: db,
	}
}

func (r *workbookRepository) AddWorkbook(ctx context.Context, operator service.OperatorInterface, param *service.WorkbookAddParameter) (*domain.WorkbookID, error) {
	_, span := tracer.Start(ctx, "workbookRepository.AddWorkbook")
	defer span.End()

	workbook := WorkbookEntity{
		Version:        1,
		CreatedBy:      operator.AppUserID().Int(),
		UpdatedBy:      operator.AppUserID().Int(),
		OrganizationID: operator.OrganizationID().Int(),
		ProblemType:    param.ProblemType,
		Name:           param.Name,
		Lang2:          param.Lang2,
		Description:    param.Description,
		Content:        param.Content,
	}
	if result := r.db.Create(&workbook); result.Error != nil {
		return nil, rsliberrors.Errorf("workbookRepository.AddWorkbook. err: %w", rslibgateway.ConvertDuplicatedError(result.Error, service.ErrWorkbookAlreadyExists))
	}

	workbookID, err := domain.NewWorkbookID(workbook.ID)
	if err != nil {
		return nil, err
	}

	return workbookID, nil
}

func (r *workbookRepository) UpdateWorkbook(ctx context.Context, operator service.OperatorInterface, workbookID *domain.WorkbookID, version int, param *service.WorkbookUpdateParameter) error {
	_, span := tracer.Start(ctx, "workbookRepository.UpdateWorkbook")
	defer span.End()

	if result := r.db.Model(&WorkbookEntity{}).
		Where("organization_id = ?", uint(operator.OrganizationID().Int())).
		Where("id = ?", workbookID.Int()).
		Where("version = ?", version).
		Updates(map[string]interface{}{
			"version":     gorm.Expr("version + 1"),
			"name":        param.Name,
			"description": param.Description,
			"content":     param.Content,
		}); result.Error != nil {
		return rsliberrors.Errorf("workbookRepository.UpdateWorkbook. err: %w", rslibgateway.ConvertDuplicatedError(result.Error, service.ErrWorkbookAlreadyExists))
	}

	return nil
}
