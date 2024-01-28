package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"strconv"

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

type EnglishSentencesOne struct {
	SrcLang2                  string `json:"srcLang2"`
	SrcAudioContent           string `json:"srcAudioContent"`
	SrcAudioLengthMillisecond int    `json:"SrcAudioLengthMillisecond"`
	SrcText                   string `json:"srcText"`
	DstLang2                  string `json:"dstLang2"`
	DstAudioContent           string `json:"dstAudioContent"`
	DstAudioLengthMillisecond int    `json:"DstAudioLengthMillisecond"`
	DstText                   string `json:"dstText"`
}

type englishSentencesCSVReader struct {
	reader1 *csv.Reader
	reader2 *csv.Reader
}

func NewEnglishSentencesCSVReader(reader1 io.Reader, reader2 io.Reader) *englishSentencesCSVReader {
	return &englishSentencesCSVReader{
		reader1: csv.NewReader(reader1),
		reader2: csv.NewReader(reader2),
	}
}

func (r *englishSentencesCSVReader) Next() (*EnglishSentencesOne, error) {
	var err error

	var text []string
	text, err = r.reader1.Read()
	if errors.Is(err, io.EOF) {
		return nil, err
	}

	var audio []string
	audio, err = r.reader2.Read()
	if errors.Is(err, io.EOF) {
		return nil, err
	}

	srcAudioLengthMillisecond, err := strconv.Atoi(audio[3])
	if err != nil {
		return nil, err
	}
	dstAudioLengthMillisecond, err := strconv.Atoi(audio[1])
	if err != nil {
		return nil, err
	}

	return &EnglishSentencesOne{
		SrcLang2:                  "ja",
		SrcAudioContent:           audio[2],
		SrcAudioLengthMillisecond: srcAudioLengthMillisecond,
		SrcText:                   text[1],
		DstLang2:                  "en",
		DstAudioContent:           audio[0],
		DstAudioLengthMillisecond: dstAudioLengthMillisecond,
		DstText:                   text[0],
	}, nil
}
