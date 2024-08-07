package network

import (
	"testing"
)

func TestNewNetMessageFromString(t *testing.T) {
	// TODO: Perform better test
	m := NewNetMessageFromString("The message")
	if m.Raw() != "The message" {
		t.Fatalf("Net Message Raw should be \"The message\" but it's %s", m.Raw())
	}
}
