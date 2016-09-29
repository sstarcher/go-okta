package okta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Client to access okta
type Client struct {
	client *http.Client
	org    string
}

// errorResponse is an error wrapper for the okta response
type errorResponse struct {
	HTTPCode int
	Response ErrorResponse
	Endpoint string
}

func (e *errorResponse) Error() string {
	return fmt.Sprintf("Error hitting api endpoint %s %s", e.Endpoint, e.Response.ErrorCode)
}

// NewClient object for calling okta
func NewClient(org string) *Client {
	return &Client{
		client: &http.Client{},
		org:    org,
	}
}

// Authenticate with okta using username and password
func (c *Client) Authenticate(username, password string) (*AuthnResponse, error) {
	var request = &AuthnRequest{
		Username: username,
		Password: password,
	}

	var response = &AuthnResponse{}
	err := c.call("authn", request, response)
	return response, err
}

// Session takes a session token and always fails
func (c *Client) Session(sessionToken string) (*SessionResponse, error) {
	var request = &SessionRequest{
		SessionToken: sessionToken,
	}

	var response = &SessionResponse{}
	err := c.call("sessions", request, response)
	return response, err
}

func (c *Client) call(endpoint string, request, response interface{}) error {
	data, _ := json.Marshal(request)

	var url = "https://" + c.org + ".okta.com/api/v1/" + endpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Accept", `application/json`)
	req.Header.Add("Content-Type", `application/json`)

	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == http.StatusOK {
		err := json.Unmarshal(body, &response)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		var errors ErrorResponse
		err = json.Unmarshal(body, &errors)

		return &errorResponse{
			HTTPCode: resp.StatusCode,
			Response: errors,
			Endpoint: url,
		}
	}

	return nil
}
