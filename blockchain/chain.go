package blockchain

import (
	"bytes"
	"encoding/json"
	"errors"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/devholic77/duckcoin/db"
	"github.com/devholic77/duckcoin/utils"
)

var b *blockChain
var blockDB *bolt.DB
var once sync.Once

var ErrNotFound = errors.New("block not found")

type blockChain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

func (b *blockChain) fromBytes(data []byte) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.Decode(b)
}

func (b *blockChain) persist() {
	db.SaveBlockChain(utils.ToBytes(b))
}

func (b *blockChain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func (b *blockChain) GetBlock(hash string) Block {
	var block Block
	block.fromBytes(db.Block(hash))
	return block
}

func BlockChain() *blockChain {
	if b == nil {
		once.Do(func() {
			b = &blockChain{"", 0}
			blockDB = db.DB()
			persistBlockChain := db.BlockChain()
			if persistBlockChain == nil {
				b.AddBlock("Genesis")
			} else {
				b.fromBytes(persistBlockChain)
			}
		})
	}
	return b
}
