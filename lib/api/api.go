package api

type AppUserInfoResponse struct {
	AppUserID      int    `json:"appUserId"`
	OrganizationID int    `json:"organizationId"`
	LoginID        string `json:"loginId"`
	Username       string `json:"username"`
}
