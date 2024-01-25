package gateway

import "encoding/json"

type EnglishSentenceModel struct {
	SrcLang2                  string `json:"srcLang2"`
	SrcAudioContent           string `json:"srcAudioContent"`
	SrcAudioLengthMillisecond int    `json:"SrcAudioLengthMillisecond"`
	SrcText                   string `json:"srcText"`
	DstLang2                  string `json:"dstLang2"`
	DstAudioContent           string `json:"dstAudioContent"`
	DstAudioLengthMillisecond int    `json:"DstAudioLengthMillisecond"`
	DstText                   string `json:"dstText"`
}

type EnglishSentencesModel struct {
	Sentences []*EnglishSentenceModel `json:"sentences"`
}

func FromEnglishSentenceModel(model *EnglishSentencesModel) ([]byte, error) {
	return json.Marshal(model)
}

func ToEnglishSentenceModel(content []byte) (*EnglishSentencesModel, error) {
	model := EnglishSentencesModel{}
	if err := json.Unmarshal(content, &model); err != nil {
		return nil, err
	}

	return &model, nil
}
