package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/devholic77/duckcoin/db"
	"github.com/devholic77/duckcoin/utils"
)

type Block struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevhash,omitempty"`
	Height   int    `json:"height"`
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func (b *Block) fromBytes(data []byte) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.Decode(b)
}

func FindBlock(hash string) *Block {
	var block Block
	utils.FromBytes(db.Block(hash), &block)
	// block.fromBytes(db.Block(hash))
	return &block
}

func createBlock(data string, prevhash string, height int) *Block {
	block := Block{
		Data:     data,
		Hash:     "",
		PrevHash: prevhash,
		Height:   height,
	}
	payload := block.Data + block.PrevHash + fmt.Sprintf("%d", block.Height)
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	block.persist()
	return &block
}

func fromBytes(data []byte) Block {
	var block Block
	decoder := json.NewDecoder(bytes.NewReader(data))
	utils.HandleErr(decoder.Decode(&block))
	return block
}
