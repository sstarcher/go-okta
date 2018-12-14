package okta

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
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
		if strings.HasPrefix(v.FactorType, "token") || v.FactorType == "push" {
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
// https://developer.okta.com/docs/api/resources/factors#verify-token-factor
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

	if resp.StatusCode != http.StatusOK {
		var errors ErrorResponse
		_ = json.Unmarshal(body, &errors)
		return nil, &errorResponse{
			HTTPCode: resp.StatusCode,
			Response: errors,
		}
	}

	var verifyResp AuthnResponse
	err = json.Unmarshal(body, &verifyResp)
	if err != nil {
		return nil, err
	}

	return &verifyResp, nil
}

// https://developer.okta.com/docs/api/resources/factors#verify-push-factor
// API diverges quite a bit from the API reference, this is the result of
// trial and error.
func (f Factor) VerifyPush(
	stateToken string,
	userAgent string,
	pollInterval time.Duration,
	pollTimeout time.Duration) (*AuthnResponse, error) {
	if f.FactorType != "push" {
		return nil, fmt.Errorf(
			"can not VerifyPush on a factor type of %s", f.FactorType)
	}

	if len(userAgent) == 0 {
		return nil, errors.New(
			"a valid HTTP User-Agent is required when verifying a push factor")
	}

	data, _ := json.Marshal(map[string]string{
		"stateToken": stateToken,
	})
	req, err := http.NewRequest("POST",
		f.Links.Verify.Href, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", userAgent)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var errors ErrorResponse
		_ = json.Unmarshal(body, &errors)
		return nil, &errorResponse{
			HTTPCode: resp.StatusCode,
			Response: errors,
		}
	}

	return pollPushResult(resp, pollInterval, time.Now().Add(pollTimeout))
}

func pollPushResult(
	resp *http.Response, interval time.Duration, until time.Time,
) (*AuthnResponse, error) {
	client := http.Client{}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var pushResult AuthnResponse
	err = json.Unmarshal(body, &pushResult)
	if err != nil {
		return nil, err
	}

	for time.Now().Before(until) &&
		pushResult.Status == "MFA_CHALLENGE" &&
		pushResult.FactorResult == "WAITING" {

		data, _ := json.Marshal(map[string]string{
			"stateToken": pushResult.StateToken,
		})

		req, err := http.NewRequest("POST",
			pushResult.Links.Next.Href, bytes.NewBuffer(data))
		if err != nil {
			return nil, err
		}

		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")

		time.Sleep(interval)

		resp, err = client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			var errors ErrorResponse
			_ = json.Unmarshal(body, &errors)
			return nil, &errorResponse{
				HTTPCode: resp.StatusCode,
				Response: errors,
			}
		}

		err = json.Unmarshal(body, &pushResult)
		if err != nil {
			return nil, err
		}
	}

	return &pushResult, nil
}
