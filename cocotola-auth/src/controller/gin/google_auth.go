package handler

type GoogleAuthParameter struct {
	OrganizationName string `json:"organizationName"`
	Code             string `json:"code"`
}
