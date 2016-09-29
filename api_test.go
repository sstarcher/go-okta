package okta

import (
	"os"
	"strings"
	"testing"
)

func TestAPIFailure(t *testing.T) {
	client := NewClient("organization", "")
	_, err := client.Authenticate("username", "password")
	if !strings.Contains(err.Error(), "E0000007") {
		t.Error("Expected E0000007, got ", err.Error())
	}
}

func TestAcceptance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	client := NewClient(os.Getenv("OKTA_ORG"), "oktapreview.com")
	_, err := client.Authenticate(os.Getenv("OKTA_USERNAME"), os.Getenv("OKTA_PASSWORD"))
	if err != nil {
		t.Error("Expected nil, got ", err.Error())
	}
}
