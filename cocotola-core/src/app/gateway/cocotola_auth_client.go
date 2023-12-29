package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/app/service"
	libapi "github.com/kujilabo/cocotola-1.21/lib/api"
	rsliberrors "github.com/kujilabo/redstart/lib/errors"
)

type cocotolaAuthClient struct {
	timeout time.Duration
}

func NewCocotolaAuthClient(timeout time.Duration) service.CocotolaAuthClient {
	return &cocotolaAuthClient{
		timeout: timeout,
	}
}

func (c *cocotolaAuthClient) RetrieveUserInfo(ctx context.Context, bearerToken string) (*libapi.AppUserInfoResponse, error) {
	client := http.Client{
		Timeout: c.timeout,
	}

	u, err := url.Parse("http://localhost:8010")
	if err != nil {
		return nil, rsliberrors.Errorf("url.Parse. err: %w", err)
	}

	u.Path = path.Join(u.Path, "v1", "auth", "userinfo")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, rsliberrors.Errorf("http.NewRequestWithContext. err: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+bearerToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, rsliberrors.Errorf("synthesize request. err: %w", err)
	}
	defer resp.Body.Close()

	response := libapi.AppUserInfoResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, rsliberrors.Errorf("json.NewDecoder. err: %w", err)
	}

	return &response, nil
}
