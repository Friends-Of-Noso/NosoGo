package data

import (
	"testing"
)

func TestNewBlockZero(t *testing.T) {
	b := NewBlockZero()
	if b.Number != 0 {
		t.Fatalf("Block number should be 0 but it's %d", b.Number)
	}
}
