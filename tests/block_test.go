package tests

import (
	"testing"

	"github.com/Friends-Of-Noso/NosoGo/data"
)

func TestNewBlockZero(t *testing.T) {
	b := data.NewBlockZero()
	if b.Number != 0 {
		t.Fatalf("Block number should be 0 but it's %d", b.Number)
	}
}
