package network

// TODO: Need to involve generics here if it makes sense...

type NetMessage struct {
	rawMessage string
}

func NewNetMessageFromString(message string) *NetMessage {
	return &NetMessage{
		rawMessage: message,
	}
}

func (nm *NetMessage) Raw() string {
	return nm.rawMessage
}
