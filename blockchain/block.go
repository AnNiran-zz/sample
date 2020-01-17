package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"time"

	log "github.com/sirupsen/logrus"
)

// Block represents a block in the node blockchain records
// Data inside block contains remote connected peer id - this is done for representative purposes
type Block struct {
	TimeStamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

// NewBlock creates a new block object, serializes it, obtain block hash and return it
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		TimeStamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
	}

	blockAsBytes := block.prepare()
	hasher := sha256.New()
	hasher.Write(blockAsBytes)
	block.Hash = hasher.Sum(nil)

	return block
}

// prepare function joins block data into a byte slice
// This function has a representative purpose; joining different parts of the
// data into the byte slice may follow a certain protocol, where each bytes and count
// will have indicative roles
// Current version is simple to have any substantial impact of holding such rules
func (b *Block) prepare() []byte {
	return bytes.Join(
		[][]byte{
			intToHex(b.TimeStamp),
			b.Data,
			b.PrevBlockHash,
		},
		[]byte{},
	)
}

// NewGenesisBlock creates and returns a genesis Block
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

// intToHex converts an int64 to a byte array
func intToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// Serialize the block
func (b *Block) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	if err := encoder.Encode(b); err != nil {
		log.Errorf("Error serializing block at %s, %s", b.TimeStamp, err.Error())
		panic(err)
	}
	return result.Bytes()
}
