package blockchain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/devholic77/duckcoin/db"
	"github.com/devholic77/duckcoin/utils"
)

type Block struct {
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevhash,omitempty"`
	Height       int    `json:"height"`
	Difficulty   int    `json:"difficulty"`
	Nonce        int    `json:"nonce"`
	TimeStamp    int    `json:"timestamp"`
	Transactions []*Tx  `json:"transaction"`
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func (b *Block) restore(data []byte) {
	utils.FromBytes(data, &b)
}

var ErrNotFound = errors.New("block not found")

func FindBlock(hash string) (*Block, error) {
	var block Block
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}

	block.restore(blockBytes)
	fmt.Println(block)
	return &block, nil
}

func createBlock(prevhash string, height int, difficulty int) *Block {
	block := Block{
		Hash:       "",
		PrevHash:   prevhash,
		Height:     height,
		Difficulty: difficulty,
		Nonce:      0,
	}
	block.mine()
	block.Transactions = Mempool.txToConfirm()
	block.persist()

	return &block
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		b.TimeStamp = int(time.Now().Unix())
		hash := utils.Hash(b)
		fmt.Printf("Target:%s\nHash:%s\nNonce:%d\n\n\n", target, hash, b.Nonce)

		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		}
		b.Nonce++
	}
}
