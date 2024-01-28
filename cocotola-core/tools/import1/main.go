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

	"github.com/kujilabo/cocotola-1.21/cocotola-core/src/domain"
	libapi "github.com/kujilabo/cocotola-1.21/lib/api"
)

func main() {
	ctx := context.Background()
	logger := slog.Default()
	bearerToken, err := authenticate(ctx, "cocotola-owner", "password")
	if err != nil {
		panic(err)
	}

	logger.Info(bearerToken)

	workbookID, err := domain.NewWorkbookID(1)
	if err != nil {
		panic(err)
	}
	workbook, err := containsWorkbook(ctx, bearerToken, workbookID)
	if err != nil {
		panic(err)
	}

	inAudioPath := "audio.csv"
	in_file, err := os.Open(inAudioPath)
	if err != nil {
		panic(err)
	}
	defer in_file.Close()
	reader := NewEnglishSentencesCSVReader(in_file)

	sentences := libapi.EnglishSentencesModel{
		Sentences: make([]*libapi.EnglishSentenceModel, 0),
	}
	for {
		logger.InfoContext(ctx, "next")
		paramOne, err := reader.Next(ctx)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			panic(err)
		}
		sentences.Sentences = append(sentences.Sentences, paramOne)
	}

	cibtebtBytes, err := json.Marshal(sentences)
	if err != nil {
		panic(err)
	}

	if workbook != nil {
		if err := updateWorkbook(ctx, bearerToken, workbookID, workbook.Version, &libapi.WorkbookUpdateParameter{
			Name:        workbook.Name,
			Description: workbook.Description,
			Content:     string(cibtebtBytes),
		}); err != nil {
			panic(err)
		}
	} else {
		if err := addWorkbook(ctx, bearerToken, &libapi.WorkbookAddParameter{
			Name:        "MY_WORKBOOK",
			ProblemType: "english_sentences",
			Lang2:       "ja",
			Description: "MY_WORKBOOK_DESCRIPTION",
			Content:     string(cibtebtBytes),
		}); err != nil {
			panic(err)
		}
	}
}

func addWorkbook(ctx context.Context, bearerToken string, param *libapi.WorkbookAddParameter) error {
	httpClient := http.Client{}

	coreEndpoint, err := url.Parse("http://localhost:8000")
	if err != nil {
		return err
	}
	u := *coreEndpoint
	u.Path = path.Join(u.Path, "v1", "workbook")

	b, err := json.Marshal(param)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(b))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+bearerToken)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(respBytes))
	}

	return nil
}

func updateWorkbook(ctx context.Context, bearerToken string, workbookID *domain.WorkbookID, version int, param *libapi.WorkbookUpdateParameter) error {
	httpClient := http.Client{}

	coreEndpoint, err := url.Parse("http://localhost:8000")
	if err != nil {
		return err
	}
	u := *coreEndpoint
	u.Path = path.Join(u.Path, "v1", "workbook", strconv.Itoa(workbookID.Int()))

	b, err := json.Marshal(param)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, u.String(), bytes.NewReader(b))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+bearerToken)

	q := req.URL.Query()
	q.Add("version", strconv.Itoa(version))
	req.URL.RawQuery = q.Encode()

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(respBytes))
	}

	return nil
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

func containsWorkbook(ctx context.Context, bearerToken string, workbookID *domain.WorkbookID) (*libapi.WorkbookRetrieveResult, error) {
	logger := slog.Default()
	httpClient := http.Client{}

	coreEndpoint, err := url.Parse("http://localhost:8000")
	if err != nil {
		return nil, err
	}
	u := *coreEndpoint
	u.Path = path.Join(u.Path, "v1", "workbook", strconv.Itoa(workbookID.Int()))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+bearerToken)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode == http.StatusOK {
		workbook := libapi.WorkbookRetrieveResult{}
		if err := json.NewDecoder(resp.Body).Decode(&workbook); err != nil {
			return nil, err
		}
		logger.InfoContext(ctx, fmt.Sprintf("%v", workbook))
		return &workbook, nil
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return nil, errors.New(string(respBytes))
}

// type EnglishSentencesOne struct {
// 	SrcLang2                  string `json:"srcLang2"`
// 	SrcAudioContent           string `json:"srcAudioContent"`
// 	SrcAudioLengthMillisecond int    `json:"SrcAudioLengthMillisecond"`
// 	SrcText                   string `json:"srcText"`
// 	DstLang2                  string `json:"dstLang2"`
// 	DstAudioContent           string `json:"dstAudioContent"`
// 	DstAudioLengthMillisecond int    `json:"DstAudioLengthMillisecond"`
// 	DstText                   string `json:"dstText"`
// }

type englishSentencesCSVReader struct {
	reader1 *csv.Reader
}

func NewEnglishSentencesCSVReader(reader1 io.Reader) *englishSentencesCSVReader {
	return &englishSentencesCSVReader{
		reader1: csv.NewReader(reader1),
	}
}

func (r *englishSentencesCSVReader) Next(ctx context.Context) (*libapi.EnglishSentenceModel, error) {
	var err error

	var audio []string
	audio, err = r.reader1.Read()
	if errors.Is(err, io.EOF) {
		return nil, err
	}

	logger := slog.Default()
	logger.InfoContext(ctx, fmt.Sprintf("length = %d", len(audio)))

	srcAudioLengthMillisecond, err := strconv.Atoi(audio[2])
	if err != nil {
		return nil, err
	}
	dstAudioLengthMillisecond, err := strconv.Atoi(audio[5])
	if err != nil {
		return nil, err
	}

	return &libapi.EnglishSentenceModel{
		SrcLang2:                  "ja",
		SrcAudioContent:           audio[4],
		SrcAudioLengthMillisecond: srcAudioLengthMillisecond,
		SrcText:                   audio[3],
		DstLang2:                  "en",
		DstAudioContent:           audio[1],
		DstAudioLengthMillisecond: dstAudioLengthMillisecond,
		DstText:                   audio[0],
	}, nil
}
