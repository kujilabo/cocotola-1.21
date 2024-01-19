package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	rsliberrors "github.com/kujilabo/redstart/lib/errors"

	"github.com/kujilabo/cocotola-1.21/cocotola-synthesizer/src/service"
	libdomain "github.com/kujilabo/cocotola-1.21/lib/domain"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type googleTTSClient struct {
	httpClient HTTPClient
	apiKey     string
}

type ttsResponse struct {
	AudioContent string `json:"audioContent"`
}

func NewGoogleTTSClient(httpClient HTTPClient, apiKey string) service.SynthesizerClient {
	return &googleTTSClient{
		httpClient: httpClient,
		apiKey:     apiKey,
	}
}

func (g *googleTTSClient) Synthesize(ctx context.Context, lang5 libdomain.Lang5, text string) (string, error) {
	ctx, span := tracer.Start(ctx, "synthesizerClient.Synthesize")
	defer span.End()

	type m map[string]interface{}

	values := m{
		"input": m{
			"text": text,
		},
		"voice": m{
			"languageCode": lang5.String(),
			"ssmlGender":   "FEMALE",
		},
		"audioConfig": m{
			"audioEncoding": "MP3",
			"speakingRate":  1,
		},
	}

	b, err := json.Marshal(values)
	if err != nil {
		return "", err
	}

	u, err := url.Parse("https://texttospeech.googleapis.com/v1/text:synthesize")
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("key", g.apiKey)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", rsliberrors.Errorf("%s", string(body))
	}

	response := ttsResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	return response.AudioContent, nil
}
