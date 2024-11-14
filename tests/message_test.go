package tests

import (
	"testing"

	"github.com/Friends-Of-Noso/NosoGo/network"
)

func TestNewNetMessageFromString(t *testing.T) {
	// TODO: Perform better test
	m := network.NewNetMessageFromString("The message")
	if m.Raw() != "The message" {
		t.Fatalf("Net Message Raw should be \"The message\" but it's %s", m.Raw())
	}
}
