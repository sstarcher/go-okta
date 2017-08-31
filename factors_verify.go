package okta

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// https://developer.okta.com/docs/api/resources/authn.html#verify-totp-factor
func (r *AuthnResponse) VerifyOTP(code string, verifyResp *AuthnResponse) (err error) {
	data, _ := json.Marshal(map[string]string{
		"passCode":   code,
		"stateToken": r.StateToken,
	})
	req, err := http.NewRequest("POST",
		r.Embedded.Factors[0].Links.Verify.Href, bytes.NewBuffer(data))
	if err != nil {
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(body, verifyResp)
	} else {
		var errors ErrorResponse
		_ = json.Unmarshal(body, &errors)
		return &errorResponse{
			HTTPCode: resp.StatusCode,
			Response: errors,
		}
	}

	return
}
