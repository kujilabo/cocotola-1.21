package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/kujilabo/cocotola-1.21/cocotola-auth/src/domain"
	rsliberrors "github.com/kujilabo/redstart/lib/errors"
)

type GoogleAuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type GoogleUserInfo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type GoogleAuthClient struct {
	client       http.Client
	clientID     string
	clientSecret string
	redirectURI  string
	grantType    string
}

func NewGoogleAuthClient(clientID, clientSecret, redirectURI string, timeout time.Duration) *GoogleAuthClient {
	return &GoogleAuthClient{
		client: http.Client{
			Timeout:   timeout,
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
		grantType:    "authorization_code",
	}
}

func (c *GoogleAuthClient) RetrieveAccessToken(ctx context.Context, code string) (*domain.AuthTokenSet, error) {
	ctx, span := tracer.Start(ctx, "googleAuthClient.RetrieveAccessToken")
	defer span.End()

	paramMap := map[string]string{
		"client_id":     c.clientID,
		"client_secret": c.clientSecret,
		"redirect_uri":  c.redirectURI,
		"grant_type":    c.grantType,
		"code":          code,
	}

	paramBytes, err := json.Marshal(paramMap)
	if err != nil {
		return nil, rsliberrors.Errorf("json.Marshal. err: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://accounts.google.com/o/oauth2/token", bytes.NewBuffer(paramBytes))
	if err != nil {
		return nil, rsliberrors.Errorf("http.NewRequestWithContext. err: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, rsliberrors.Errorf("failed to retrieve access token.err: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, rsliberrors.Errorf("io.ReadAll. err: %w", err)
		}
		// logger.Debugf("status:%d", resp.StatusCode)
		// logger.Debugf("Resp:%s", string(respBytes))
		return nil, errors.New(string(respBytes))
	}

	googleAuthResponse := GoogleAuthResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&googleAuthResponse); err != nil {
		return nil, rsliberrors.Errorf("json.NewDecoder. err: %w", err)
	}
	// logger.Infof("RetrieveAccessToken:%s", googleAuthResponse.AccessToken)

	return &domain.AuthTokenSet{
		AccessToken:  googleAuthResponse.AccessToken,
		RefreshToken: googleAuthResponse.RefreshToken,
	}, nil
}

func (c *GoogleAuthClient) RetrieveUserInfo(ctx context.Context, googleAuthResponse *domain.AuthTokenSet) (*domain.UserInfo, error) {
	ctx, span := tracer.Start(ctx, "googleAuthClient.RetrieveUserInfo")
	defer span.End()

	// logger := log.FromContext(ctx)

	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v1/userinfo", http.NoBody)
	if err != nil {
		return nil, rsliberrors.Errorf("http.NewRequestWithContext. err: %w", err)
	}

	q := req.URL.Query()
	q.Add("alt", "json")
	q.Add("access_token", googleAuthResponse.AccessToken)
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, rsliberrors.Errorf("c.client.Do. err: %w", err)
	}
	defer resp.Body.Close()

	// logger.Debugf("access_token:%s", googleAuthResponse.AccessToken)
	// logger.Debugf("status:%d", resp.StatusCode)

	googleUserInfo := GoogleUserInfo{}
	if err := json.NewDecoder(resp.Body).Decode(&googleUserInfo); err != nil {
		return nil, rsliberrors.Errorf("json.NewDecoder. err: %w", err)
	}

	return &domain.UserInfo{
		Email: googleUserInfo.Email,
		Name:  googleUserInfo.Name,
	}, nil
}
