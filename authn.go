package okta

import (
	"time"
)

type ErrorResponse struct {
	ErrorCode    string `json:"errorCode"`
	ErrorSummary string `json:"errorSummary"`
	ErrorLink    string `json:"errorLink"`
	ErrorID      string `json:"errorId"`
	ErrorCauses  []struct {
		ErrorSummary string `json:"errorSummary"`
	} `json:"errorCauses"`
}

type AuthnRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RelayState string `json:"relayState"`
	Options    struct {
		MultiOptionalFactorEnroll bool `json:"multiOptionalFactorEnroll"`
		WarnBeforePasswordExpired bool `json:"warnBeforePasswordExpired"`
	} `json:"options"`
}

type AuthnResponse struct {
	StateToken   string    `json:"stateToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
	Status       string    `json:"status"`
	RelayState   string    `json:"relayState"`
	FactorResult string    `json:"factorResult"`
	SessionToken string    `json:"sessionToken"`
	Embedded     struct {
		User struct {
			ID              string    `json:"id"`
			PasswordChanged time.Time `json:"passwordChanged"`
			Profile         struct {
				Login     string `json:"login"`
				FirstName string `json:"firstName"`
				LastName  string `json:"lastName"`
				Locale    string `json:"locale"`
				TimeZone  string `json:"timeZone"`
			} `json:"profile"`
		} `json:"user"`
		Factors []Factor `json:"factors"`
		Policy  struct {
			AllowRememberDevice             bool `json:"allowRememberDevice"`
			RememberDeviceLifetimeInMinutes int  `json:"rememberDeviceLifetimeInMinutes"`
			RememberDeviceByDefault         bool `json:"rememberDeviceByDefault"`
		} `json:"policy"`
	} `json:"_embedded"`
	Links struct {
		Cancel struct {
			Href  string `json:"href"`
			Hints struct {
				Allow []string `json:"allow"`
			} `json:"hints"`
		} `json:"cancel"`
		Next struct {
			Name  string `json:"name"`
			Href  string `json:"href"`
			Hints struct {
				Allow []string `json:"allow"`
			} `json:"hints"`
		}
	} `json:"_links"`
}
