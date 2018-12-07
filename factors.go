package okta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Factor struct {
	ID         string `json:"id"`
	FactorType string `json:"factorType"`
	Provider   string `json:"provider"`
	VendorName string `json:"vendorName"`
	Profile    struct {
		CredentialID string `json:"credentialId"`
	} `json:"profile"`
	Links struct {
		Verify struct {
			Href  string `json:"href"`
			Hints struct {
				Allow []string `json:"allow"`
			} `json:"hints"`
		} `json:"verify"`
	} `json:"_links"`
}

func (r *AuthnResponse) GetSupportedFactors() []Factor {
	var supported []Factor

	for _, v := range r.Embedded.Factors {
		postAllowed := false
		if strings.HasPrefix(v.FactorType, "token") {
			for _, verb := range v.Links.Verify.Hints.Allow {
				if verb == "POST" {
					postAllowed = true
				}
			}

			if postAllowed {
				supported = append(supported, v)
			}
		}
	}

	return supported
}

// https://developer.okta.com/docs/api/resources/factors#verify-totp-factor
func (f Factor) VerifyOTP(stateToken string, code string) (*AuthnResponse, error) {
	if !strings.HasPrefix(f.FactorType, "token") {
		return nil, fmt.Errorf(
			"can not VerifyOTP on a factor type of %s", f.FactorType)
	}

	data, _ := json.Marshal(map[string]string{
		"passCode":   code,
		"stateToken": stateToken,
	})
	req, err := http.NewRequest("POST",
		f.Links.Verify.Href, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var verifyResp AuthnResponse
	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(body, &verifyResp)
		if err != nil {
			return nil, err
		}
	} else {
		var errors ErrorResponse
		_ = json.Unmarshal(body, &errors)
		return nil, &errorResponse{
			HTTPCode: resp.StatusCode,
			Response: errors,
		}
	}

	return &verifyResp, nil
}
