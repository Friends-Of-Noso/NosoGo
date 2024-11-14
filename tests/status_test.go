package tests

import (
	"testing"

	"github.com/Friends-Of-Noso/NosoGo/core"
)

func TestNewNetworkStatus(t *testing.T) {
	ns := core.NewNetworkStatus()
	if ns.LastBlock != 0 {
		t.Fatalf("Network Status LastBlock should be 0 but it's %d", ns.LastBlock)
	}
}
