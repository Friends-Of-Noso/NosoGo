package node

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

const (
	BLOCKS_SUB      = "blocks"
	CONNECTIONS_SUB = "connections"
)

type (
	PubSubTopics       map[string]*pubsub.Topic
	PubSubSubscription map[string]*pubsub.Subscription
)
