package crypt

import (
	"crypto/rand"
	"io"
	mrand "math/rand"

	crypto "github.com/libp2p/go-libp2p-crypto"
	log "github.com/sirupsen/logrus"
)

var r io.Reader

// GenerateKey generates private key using RSA cryptosystem
func GenerateKey(srcPort int, debug bool) crypto.PrivKey {
	// If debug is enabled, source port is used for generating peer ID
	// this is only useful for debugging purposes
	// this will always generate the same node ID for multiple executions
	if debug {
		r = mrand.New(mrand.NewSource(int64(srcPort)))
	} else {
		// use random reader
		r = rand.Reader
	}

	privKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		log.Errorf("Error generating private key: %s", err.Error())
		panic(err)
	}

	return privKey
}
