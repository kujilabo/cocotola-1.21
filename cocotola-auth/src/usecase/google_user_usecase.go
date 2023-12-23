package usecase

import (
	"context"

	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/service"
	rsliberrors "github.com/kujilabo/redstart/lib/errors"
)

type TokenSet struct {
	AccessToken  string
	RefreshToken string
}
type GoogleAuthClient interface {
	RetrieveAccessToken(ctx context.Context, code string) (*GoogleAuthResponse, error)
	RetrieveUserInfo(ctx context.Context, googleAuthResponse *GoogleAuthResponse) (*GoogleUserInfo, error)
}

type GoogleAuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type GoogleUserInfo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type GoogleUserUsecaseInterface interface {
	RetrieveAccessToken(ctx context.Context, code string) (*GoogleAuthResponse, error)

	// RetrieveUserInfo(ctx context.Context, GoogleAuthResponse *GoogleAuthResponse) (*GoogleUserInfo, error)

	// RegisterAppUser(ctx context.Context, googleUserInfo *GoogleUserInfo, googleAuthResponse *GoogleAuthResponse, organizationName string) (*TokenSet, error)
}

type GoogleUserUsecase struct {
	transactionManager service.TransactionManager
	googleAuthClient   GoogleAuthClient
}

func NewGoogleUserUsecase(transactionManager service.TransactionManager, googleAuthClient GoogleAuthClient) *GoogleUserUsecase {
	return &GoogleUserUsecase{
		transactionManager: transactionManager,
		googleAuthClient:   googleAuthClient,
	}
}

func (s *GoogleUserUsecase) RetrieveAccessToken(ctx context.Context, code string) (*GoogleAuthResponse, error) {
	resp, err := s.googleAuthClient.RetrieveAccessToken(ctx, code)
	if err != nil {
		return nil, rsliberrors.Errorf(". err: %w", err)
	}

	return resp, nil
}

func (s *GoogleUserUsecase) RetrieveUserInfo(ctx context.Context, googleAuthResponse *GoogleAuthResponse) (*GoogleUserInfo, error) {
	info, err := s.googleAuthClient.RetrieveUserInfo(ctx, googleAuthResponse)
	if err != nil {
		return nil, rsliberrors.Errorf(". err: %w", err)
	}

	return info, nil
}
