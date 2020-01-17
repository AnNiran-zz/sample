package blockchain

import (
	"sample/database"
	"sync"

	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
)

// Blockchain struct contains database connection and a tip of the current blockchain
// in this case the blockchain is really a chain, instead of a tree to select a best state from
type Blockchain struct {
	tip []byte
	db  *bolt.DB
	sync.RWMutex
}

// dbFile would have the same filename across node directories
const dbFile = "blockchain.db"

// NewBlockchain opens a blockchain or create new and creates a genesis file
func NewBlockchain(nodepath string) *Blockchain {
	var tip []byte
	db := database.OpenDB(nodepath, dbFile)

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(database.BucketName))

		if b == nil {
			log.Info("No existing blockchain found. Creating a new one.")
			genesis := NewGenesisBlock()

			// create bucket with the standard name
			// since all objects are identical - we create bucket with same name for each blockchain
			b, err := tx.CreateBucket([]byte(database.BucketName))
			if err != nil {
				log.Errorf("Error creating new bucket: %s", err.Error())
				panic(err)
			}

			// set bucket value key
			if err = b.Put(genesis.Hash, genesis.Serialize()); err != nil {
				log.Errorf("Error setting the bucket value key: %s", err.Error())
				panic(err)
			}

			if err = b.Put([]byte("l"), genesis.Hash); err != nil {
				log.Error(err.Error())
				panic(err)
			}
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	if err != nil {
		log.Error(err.Error())
		panic(err)
	}

	return &Blockchain{tip: tip, db: db}
}

// AddBlock saves provided data in a block in the blockchain
// function is save for concurrently writing data
//
// Concurency save can be used if the functionality is extended, in the present sample it is not necessary
func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	bc.Lock()
	err := bc.db.View(func(tx *bolt.Tx) error {
		// obtain bucket
		b := tx.Bucket([]byte(database.BucketName))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Error(err.Error())
		panic(err)
	}

	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(database.BucketName))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Error(err.Error())
			panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Error(err.Error())
			panic(err)
		}

		bc.tip = newBlock.Hash

		return nil
	})

	bc.Unlock()
	log.Infof("Saved peerID: %s into connections ledger blockhain", data)
}
