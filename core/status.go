package core

type NetworkStatus struct {
	LastBlock uint64
}

func NewNetworkStatus() *NetworkStatus {
	return &NetworkStatus{
		LastBlock: 0,
	}
}
