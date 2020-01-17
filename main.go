package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sample/blockchain"
	"sample/crypt"
	"sync"
	"time"

	"github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"

	libp2p "github.com/libp2p/go-libp2p"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	protocol "github.com/libp2p/go-libp2p-protocol"
)

var bc *blockchain.Blockchain
var wg sync.WaitGroup

func main() {
	// Receive command line arguments for port and debug boolean value
	srcPort := flag.Int("sp", 0, "Port number")
	debug := flag.Bool("debug", false, "Generate the same node ID on each execution")
	flag.Parse()

	// Generate RSA 2048-bit key using the provided port number and debug option
	privKey := crypt.GenerateKey(*srcPort, *debug)

	// Create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up multiaddress cross-protocol format
	srcMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *srcPort))
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}

	// Create node using the set up private key, multiaddress and set up replay enabling
	node, err := libp2p.New(
		ctx,
		libp2p.ListenAddrs(srcMultiAddr),
		libp2p.Identity(privKey),
		libp2p.EnableRelay(),
		libp2p.ConnectionManager(
			connmgr.NewConnManager(
				110,
				550,
				time.Minute,
			),
		),
	)
	if err != nil {
		log.Errorf("Error creating node: %s", err.Error())
		panic(err)
	}
	log.Info("Node created with ID:", node.ID())
	log.Info("Node listen address is:", node.Addrs())

	// Fill blockchain instance
	bc = connBlockchain(node.ID().String())

	// Subscribe to topic1 and topic2
	setUpPubSub(ctx, node)

	node.SetStreamHandler(protocol.ID("service/1.0"), handleService0)

	// Start a DHT that will be used in peer discovery
	// Each peer will keep a local copy of the DHT
	kadDHT, err := dht.New(ctx, node)
	if err != nil {
		log.Errorf("Error creating a dht local copy: %s", err.Error())
		panic(err)
	}

	if err = kadDHT.Bootstrap(ctx); err != nil {
		log.Error(err.Error())
		panic(err)
	}

	// Connect to bootstrap nodes first
	bootStrapPeers := dht.DefaultBootstrapPeers

	for _, peerAdr := range bootStrapPeers {
		peerInfo, _ := peer.AddrInfoFromP2pAddr(peerAdr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := node.Connect(ctx, *peerInfo); err != nil {
				log.Error(err.Error())

			} else {
				log.Info("Connection established with bootstrap node:", *peerInfo)
				bc.AddBlock(peerInfo.String())
			}
		}()
	}
	wg.Wait()

	// Announce node
	log.Info("Announcing local node...")
	routeDiscovery := discovery.NewRoutingDiscovery(kadDHT)
	discovery.Advertise(ctx, routeDiscovery, "location")
	log.Info("Node successfully announced")

	// Seek for other nodes that announced themselves
	log.Info("Search for to other peers...")
	peerChn, err := routeDiscovery.FindPeers(ctx, "location")
	if err != nil {
		panic(err)
	}

	for peer := range peerChn {
		if peer.ID == node.ID() {
			continue
		}
		log.Info("Found peer:", peer)
		log.Info("Connecting to peer:", peer)
		_, err = node.NewStream(ctx, peer.ID, protocol.ID("service/1.0"))

		if err != nil {
			log.Errorf("Connection failed: %s", err.Error())
			continue
		} else {
			// set up libp2p-pubsub
			setUpPubSub(ctx, node)
		}

		log.Info("Connected to:", peer)
		bc.AddBlock(peer.ID.String())
	}

	select {}
}

// Check if corresponding folders exists for using connections ledger/blockchain
func connBlockchain(peerID string) *blockchain.Blockchain {
	path := filepath.Join(os.Getenv("HOME"), ".boltdb", peerID)
	if _, err := os.Stat(path); err != nil {
		os.MkdirAll(path, 0755)
	}

	return blockchain.NewBlockchain(path)
}
