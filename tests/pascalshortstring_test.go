package tests

import (
	"testing"

	"github.com/Friends-Of-Noso/NosoGo/utils"
)

func TestNewPascalShortStringEmpty(t *testing.T) {
	nps := utils.NewPascalShortString("", 10)
	if nps.String() != "" {
		t.Fatalf("Content is not an empty string but is %s", nps.String())
	}
	if nps.Len() != 0 {
		t.Fatalf("Length is not 0 but is %d", nps.Len())
	}
	if nps.Capacity() != 10 {
		t.Fatalf("Capacity is not 10 but is %d", nps.Capacity())
	}
}

func TestNewPascalShortString(t *testing.T) {
	nps := utils.NewPascalShortString("Hello World!", 15)
	if nps.String() != "Hello World!" {
		t.Fatalf("Content is not \"Hello World!\" but is %s", nps.String())
	}
	if nps.Len() != 12 {
		t.Fatalf("Length is not 12 but is %d", nps.Len())
	}
	if nps.Capacity() != 15 {
		t.Fatalf("Capacity is not 10 but is %d", nps.Capacity())
	}
}

// TODO: Add test for new content
func TestPascalShortStringNewContent(t *testing.T) {
	pss := utils.NewPascalShortString("abc", 5)
	pss.NewContent("aaaaaa")
	if pss.String() != "aaaaa" {
		t.Fatalf("Content is not \"aaaaa\" but %s", pss.String())
	}
	pss.NewContent("bbb")
	if string(pss.Raw()) != "bbbaa" {
		t.Fatalf("Content is not \"bbbaa\" but %s", string(pss.Raw()))
	}
}
