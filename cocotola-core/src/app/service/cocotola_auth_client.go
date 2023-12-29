package service

import (
	"context"

	libapi "github.com/kujilabo/cocotola-1.21/lib/api"
)

type CocotolaAuthClient interface {
	RetrieveUserInfo(ctx context.Context, bearerToken string) (*libapi.AppUserInfoResponse, error)
}
