package service

import (
	"context"
	"errors"

	libdomain "github.com/kujilabo/cocotola-1.21/lib/domain"

	"github.com/kujilabo/cocotola-1.21/cocotola-synthesizer/src/domain"
)

var ErrAudioNotFound = errors.New("Audio not found")

type AudioRepository interface {
	AddAudio(ctx context.Context, lang5 *libdomain.Lang5, text, audioContent string) (*domain.AudioID, error)

	FindAudioByAudioID(ctx context.Context, audioID *domain.AudioID) (*domain.AudioModel, error)

	FindByLangAndText(ctx context.Context, lang5 *libdomain.Lang5, text string) (*domain.AudioModel, error)

	FindAudioIDByText(ctx context.Context, lang5 *libdomain.Lang5, text string) (*domain.AudioID, error)
}
