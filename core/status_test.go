package core

import (
	"testing"
)

func TestNewNetworkStatus(t *testing.T) {
	ns := NewNetworkStatus()
	if ns.LastBlock != 0 {
		t.Fatalf("Network Status LastBlock should be 0 but it's %d", ns.LastBlock)
	}
}
