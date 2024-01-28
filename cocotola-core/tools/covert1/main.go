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
	"os"
	"path"
	"strconv"
	"time"

	libapi "github.com/kujilabo/cocotola-1.21/lib/api"
)

func main() {
	ctx := context.Background()
	inTextPath := "text.csv"
	outAudioPath := "audio.csv"
	in_file, err := os.Open(inTextPath)
	if err != nil {
		panic(err)
	}
	defer in_file.Close()

	out_file, err := os.Create(outAudioPath)
	if err != nil {
		panic(err)
	}
	defer out_file.Close()
	synthesizerClient := synthesizeClient{
		endpoint: "http://localhost:8020",
		username: "username",
		password: "password",
	}
	reader := NewEnglishSentencesCSVReader(in_file, out_file, &synthesizerClient)

	logger := slog.Default()
	for {
		logger.InfoContext(ctx, "next")
		err := reader.Next(ctx)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			panic(err)
		}
	}
}

type englishSentencesCSVReader struct {
	reader1          *csv.Reader
	writer1          *csv.Writer
	synthesizeClient *synthesizeClient
}

func NewEnglishSentencesCSVReader(reader1 io.Reader, writer1 io.Writer, synthesizeClient *synthesizeClient) *englishSentencesCSVReader {
	return &englishSentencesCSVReader{
		reader1:          csv.NewReader(reader1),
		writer1:          csv.NewWriter(writer1),
		synthesizeClient: synthesizeClient,
	}
}

func (r *englishSentencesCSVReader) Next(ctx context.Context) error {
	var err error

	var text []string
	text, err = r.reader1.Read()
	if errors.Is(err, io.EOF) {
		return err
	}

	srcText := text[1]
	dstText := text[0]

	srcAudioContent, srcAudioLength, err := r.synthesizeClient.synthesize(ctx, "ja-JP", "ja-JP-Neural2-C", srcText)
	if err != nil {
		return err
	}
	dstAudioContent, dstAudioLength, err := r.synthesizeClient.synthesize(ctx, "en-US", "en-US-Neural2-A", dstText)
	if err != nil {
		return err
	}

	if err := r.writer1.Write([]string{srcText, srcAudioContent, strconv.Itoa(int(srcAudioLength.Milliseconds())), dstText, dstAudioContent, strconv.Itoa(int(dstAudioLength.Milliseconds()))}); err != nil {
		return err
	}

	r.writer1.Flush()

	return nil
}

type synthesizeClient struct {
	endpoint string
	username string
	password string
}

func (c *synthesizeClient) synthesize(ctx context.Context, lang5, voice, text string) (string, time.Duration, error) {
	logger := slog.Default()
	logger.InfoContext(ctx, "synthesize")
	httpClient := http.Client{}

	authEndpoint, err := url.Parse(c.endpoint)
	if err != nil {
		return "", 0, err
	}
	u := *authEndpoint
	u.Path = path.Join(u.Path, "v1", "synthesize", "synthesize")

	param := libapi.SynthesizeParameter{
		Lang5: lang5,
		Voice: voice,
		Text:  text,
	}
	b, err := json.Marshal(param)
	if err != nil {
		return "", 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(b))
	if err != nil {
		return "", 0, err
	}
	req.SetBasicAuth(c.username, c.password)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", 0, err
		}
		return "", 0, fmt.Errorf("%d, %s", resp.StatusCode, string(respBytes))
	}
	synthesizeResponse := libapi.SynthesizeResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&synthesizeResponse); err != nil {
		return "", 0, err
	}

	return synthesizeResponse.AudioContent, time.Duration(synthesizeResponse.AudioLengthMillisecond) * time.Millisecond, nil
}
