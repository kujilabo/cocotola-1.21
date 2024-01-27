package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"

	libapi "github.com/kujilabo/cocotola-1.21/lib/api"
)

func main() {
	ctx := context.Background()
	logger := slog.Default()
	bearerToken, err := authenticate(ctx, "cocotola-owner", "password")
	if err != nil {
		panic(err)
	}
	contains, err := containsWorkbook(ctx, bearerToken)
	if err != nil {
		panic(err)
	}
	logger.InfoContext(ctx, fmt.Sprintf("%v", contains))
}

func authenticate(ctx context.Context, loginID, password string) (string, error) {
	httpClient := http.Client{}

	authEndpoint, err := url.Parse("http://localhost:8010")
	if err != nil {
		return "", err
	}
	u := *authEndpoint
	u.Path = path.Join(u.Path, "v1", "password", "authenticate")

	param := libapi.PasswordAuthParameter{
		LoginID:          loginID,
		Password:         password,
		OrganizationName: "cocotola",
	}
	b, err := json.Marshal(param)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return "", errors.New(string(respBytes))
	}
	authTokenSet := libapi.AuthResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&authTokenSet); err != nil {
		return "", err
	}

	return *authTokenSet.AccessToken, nil
}

func containsWorkbook(ctx context.Context, bearerToken string) (bool, error) {
	httpClient := http.Client{}

	coreEndpoint, err := url.Parse("http://localhost:8000")
	if err != nil {
		return false, err
	}
	u := *coreEndpoint
	u.Path = path.Join(u.Path, "v1", "workbook", "1")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "Bearer "+bearerToken)

	resp, err := httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	return false, errors.New(string(respBytes))
}
