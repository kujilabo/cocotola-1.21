package usecase

import (
	"context"
	"errors"

	libdomain "github.com/kujilabo/cocotola-1.21/lib/domain"

	rsliberrors "github.com/kujilabo/redstart/lib/errors"

	"github.com/kujilabo/cocotola-1.21/cocotola-synthesizer/src/domain"
	"github.com/kujilabo/cocotola-1.21/cocotola-synthesizer/src/service"
)

type SynthesizerUsecase struct {
	txManager         service.TransactionManager
	nonTxManager      service.TransactionManager
	synthesizerClient service.SynthesizerClient
}

func NewSynthesizerUsecase(txManager service.TransactionManager, nonTxManager service.TransactionManager, synthesizerClient service.SynthesizerClient) *SynthesizerUsecase {
	return &SynthesizerUsecase{

		txManager:         txManager,
		nonTxManager:      nonTxManager,
		synthesizerClient: synthesizerClient,
	}
}

func (u *SynthesizerUsecase) Synthesize(ctx context.Context, lang5 *libdomain.Lang5, text string) (*domain.AudioModel, error) {
	var audioModel *domain.AudioModel

	if err := u.txManager.Do(ctx, func(rf service.RepositoryFactory) error {
		// try to find the audio content from the DB
		repo := rf.NewAudioRepository(ctx)
		if tmpAudioModel, err := repo.FindByLangAndText(ctx, lang5, text); err == nil {
			audioModel = tmpAudioModel
			return nil
		} else if !errors.Is(err, service.ErrAudioNotFound) {
			return rsliberrors.Errorf("FindByLangAndText. err: %w", err)
		}
		return service.ErrAudioNotFound
	}); err != nil {
		return nil, err
	}

	var audioID *domain.AudioID
	if err := u.txManager.Do(ctx, func(rf service.RepositoryFactory) error {
		// synthesize text via the Web API
		repo := rf.NewAudioRepository(ctx)
		audioContent, err := u.synthesizerClient.Synthesize(ctx, lang5, "FEMALE", text)
		if err != nil {
			return rsliberrors.Errorf("to u.synthesizerClient.Synthesize. err: %w", err)
		}
		tmpAudioID, err := repo.AddAudio(ctx, lang5, text, audioContent)
		if err != nil {
			return rsliberrors.Errorf("repo.AddAudio. err: %w", err)
		}
		audioID = tmpAudioID
		return nil
	}); err != nil {
		return nil, err
	}

	if err := u.txManager.Do(ctx, func(rf service.RepositoryFactory) error {
		// try to find the audio content from the DB
		repo := rf.NewAudioRepository(ctx)
		tmpAudioModel, err := repo.FindAudioByAudioID(ctx, audioID)
		if err != nil {
			return rsliberrors.Errorf("repo.FindAudioByAudioID. err: %w", err)
		}
		audioModel = tmpAudioModel
		return nil
	}); err != nil {
		return nil, err
	}

	return audioModel, nil
}

func (u *SynthesizerUsecase) FindAudioByID(ctx context.Context, audioID *domain.AudioID) (*domain.AudioModel, error) {

	var audio *domain.AudioModel
	if err := u.nonTxManager.Do(ctx, func(rf service.RepositoryFactory) error {
		repo := rf.NewAudioRepository(ctx)
		tmpAudio, err := repo.FindAudioByAudioID(ctx, audioID)
		if err != nil {
			return err
		}

		audio = tmpAudio
		return nil
	}); err != nil {
		return nil, err
	}

	return audio, nil
}
