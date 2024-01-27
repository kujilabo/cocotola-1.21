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
