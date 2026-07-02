package cli

import (
	"strings"
	"testing"
)

// Read commands need a live client, so offline tests assert arg validation:
// a registered get command with no <id> must fail with cobra's arg error,
// not "unknown command".
func TestGetCommandsRequireID(t *testing.T) {
	cases := [][]string{
		{"firewall", "get"},
		{"devicecontrol", "get"},
		{"remoteops", "get"},
		{"updates", "get"},
	}
	for _, args := range cases {
		_, err := runCLI(t, args...)
		if err == nil {
			t.Fatalf("%v: expected missing-argument error", args)
		}
		if strings.Contains(err.Error(), "unknown command") {
			t.Fatalf("%v: command not registered: %v", args, err)
		}
	}
}
