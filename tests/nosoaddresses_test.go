package tests

import (
	"testing"

	"github.com/Friends-Of-Noso/NosoGo/legacy"
)

const (
	pubKey   = "BDAXq+mZYkwN5DS7ABR5VruS1u1ZMkiLKip8IHWjJJ4YP3bDgK45Ey13dpijXsNWdOaTeSOO1jlCEo3OxftQel8="
	addressN = "NuxYnPPYEqFMw3UM8j3hLppXsF8dEk"
	addressM = "MuxYnPPYEqFMw3UM8j3hLppXsF8dEk"
)

func TestGetAddressFromPublicKeyN(t *testing.T) {
	address := legacy.GetAddressFromPublicKey(pubKey, 0)
	if address != addressN {
		t.Errorf("Error: Expected %s Got %s", addressN, address)
	}
}

func TestGetAddressFromPublicKeyM(t *testing.T) {
	address := legacy.GetAddressFromPublicKey(pubKey, 1)
	if address != addressM {
		t.Errorf("Error: Expected %s Got %s", addressM, address)
	}
}
