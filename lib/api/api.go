package api

type AppUserInfoResponse struct {
	AppUserID      int    `json:"appUserId"`
	OrganizationID int    `json:"organizationId"`
	LoginID        string `json:"loginId"`
	Username       string `json:"username"`
}

type PasswordAuthParameter struct {
	LoginID          string `json:"loginId"`
	Password         string `json:"password"`
	OrganizationName string `json:"organizationName"`
}

type AuthResponse struct {
	AccessToken  *string `json:"accessToken"`
	RefreshToken *string `json:"refreshToken"`
}

type RefreshTokenParameter struct {
	RefreshToken string `json:"refreshToken"`
}

type SynthesizeParameter struct {
	Lang5 string `json:"lang5" binding:"required,len=5"`
	Voice string `json:"voice"`
	Text  string `json:"text"`
}

type SynthesizeResponse struct {
	AudioContent           string `json:"audioContent"`
	AudioLengthMillisecond int    `json:"audioLengthMillisecond"`
}
