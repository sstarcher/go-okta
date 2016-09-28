package main

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
	ExpiresAt    time.Time `json:"expiresAt"`
	Status       string    `json:"status"`
	RelayState   string    `json:"relayState"`
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
	} `json:"_embedded"`
}

type SessionRequest struct {
	SessionToken string `json:"sessionToken"`
}

type SessionResponse struct {
	ID                       string      `json:"id"`
	Login                    string      `json:"login"`
	UserID                   string      `json:"userId"`
	ExpiresAt                time.Time   `json:"expiresAt"`
	Status                   string      `json:"status"`
	LastPasswordVerification time.Time   `json:"lastPasswordVerification"`
	LastFactorVerification   interface{} `json:"lastFactorVerification"`
	Amr                      []string    `json:"amr"`
	Idp                      struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"idp"`
	MfaActive bool `json:"mfaActive"`
	Links     struct {
		Self struct {
			Href  string `json:"href"`
			Hints struct {
				Allow []string `json:"allow"`
			} `json:"hints"`
		} `json:"self"`
		Refresh struct {
			Href  string `json:"href"`
			Hints struct {
				Allow []string `json:"allow"`
			} `json:"hints"`
		} `json:"refresh"`
		User struct {
			Name  string `json:"name"`
			Href  string `json:"href"`
			Hints struct {
				Allow []string `json:"allow"`
			} `json:"hints"`
		} `json:"user"`
	} `json:"_links"`
}
