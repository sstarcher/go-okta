package okta

import (
	"strings"
	"testing"
)

func TestAPIFailure(t *testing.T) {
	client := NewClient("organization")
	_, err := client.Authenticate("username", "password")
	if !strings.Contains(err.Error(), "E0000007") {
		t.Error("Expected E0000007, got ", err.Error())
	}
}
