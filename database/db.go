package database

import (
	"path/filepath"

	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
)

const (
	// BucketName is the same across all nodes records
	BucketName     = "service/1.0"
)

// Error strings
const (
	errDbNotOpenStr = "database is not open"
	errTxClosedStr  = "database tx is closed"
)

// OpenDB creates or open existing db file
func OpenDB(path, file string) *bolt.DB {
	db, err := bolt.Open(filepath.Join(path, file), 0600, nil)
	if err != nil {
		log.Errorf("Error creating database: %s", err.Error())
		panic(err)
	}

	return db
}
