package main

import (
	"context"

	"github.com/libp2p/go-libp2p-core/host"
	net "github.com/libp2p/go-libp2p-net"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
)

var topic1 = "topic1subscription"
var topic2 = "topic2subscription"

func handleService0(stream net.Stream) {
	log.Info("Connection established with %s", stream.Conn().RemotePeer().String())

	//
}

func setUpPubSub(ctx context.Context, node host.Host) {
	// Set up pubsub flood
	fldSubRouter, err := pubsub.NewFloodSub(ctx, node)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Subscribe to topic1
	_, err = fldSubRouter.Join(topic1)
	if err != nil {
		log.Errorf("Error subscribing to %s: %s", topic1, err.Error())
	}
	log.Infof("Node %s subscribed to %s", node.ID().String(), topic1)

	// Subscribe to topic2
	_, err = fldSubRouter.Join(topic2)
	if err != nil {
		log.Errorf("Error subscribing to %s: %s", topic2, err.Error())
	}
	log.Infof("Node %s subscribed to %s", node.ID().String(), topic2)

	// ...

}
