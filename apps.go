package okta

type AppLinks []struct {
	AppAssignmentID  string `json:"appAssignmentId"`
	AppInstanceID    string `json:"appInstanceId"`
	AppName          string `json:"appName"`
	CredentialsSetup bool   `json:"credentialsSetup"`
	Hidden           bool   `json:"hidden"`
	ID               string `json:"id"`
	Label            string `json:"label"`
	LinkURL          string `json:"linkUrl"`
	LogoURL          string `json:"logoUrl"`
	SortOrder        int64  `json:"sortOrder"`
}
