package blockchain

import (
	"fmt"
	"sync"

	"github.com/devholic77/duckcoin/db"
	"github.com/devholic77/duckcoin/utils"
)

var b *blockChain
var once sync.Once

const (
	defaultDifficulty  = 2
	difficultyInterval = 5
	blockInterval      = 2
	allowedRange       = 2
)

type blockChain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
}

func (b *blockChain) restore(data []byte) {
	utils.FromBytes(data, &b)
}

func persist(b *blockChain) {
	db.SaveBlockChain(utils.ToBytes(b))
}

func Blocks(b *blockChain) []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

func Txs(b *blockChain) []*Tx {
	var txs []*Tx
	for _, block := range Blocks(b) {
		txs = append(txs, block.Transactions...)
	}
	return txs
}

func FindTx(b *blockChain, targetID string) *Tx {
	for _, tx := range Txs(b) {
		if tx.Id == targetID {
			return tx
		}
	}
	return nil
}

func (b *blockChain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height+1, getDifficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	persist(b)
}

func (b *blockChain) GetBlock(hash string) Block {
	var block Block
	block.restore(db.Block(hash))
	return block
}

func recalculateDifficulty(b *blockChain) int {
	allBlocks := Blocks(b)
	neweshBlock := allBlocks[0]
	lastRecalcurateBlock := allBlocks[difficultyInterval-1]

	acutalTime := (neweshBlock.TimeStamp / 60) - (lastRecalcurateBlock.TimeStamp / 60)
	expectedTime := difficultyInterval * blockInterval

	// minute buffer +-2
	if acutalTime <= (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if acutalTime >= (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}

func getDifficulty(b *blockChain) int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		fmt.Println("called recaculate")
		return recalculateDifficulty(b)
	}
	return b.CurrentDifficulty
}

func UTxOutsByAddress(address string, b *blockChain) []*UTxOut {
	var uTxOuts []*UTxOut
	creatorTxs := make(map[string]bool)
	for _, block := range Blocks(b) {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Signature == "COINBASE" {
					break
				}

				if FindTx(b, input.TxID).TxOuts[input.Index].Address == address {
					creatorTxs[input.TxID] = true
				}
			}

			for index, output := range tx.TxOuts {
				if _, ok := creatorTxs[tx.Id]; !ok {
					if output.Address == address {
						uTxOut := &UTxOut{
							tx.Id,
							index,
							output.Amount,
						}
						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}
	return uTxOuts
}

func BalanceByAddress(address string, b *blockChain) int {
	txOuts := UTxOutsByAddress(address, b)
	var amount int
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
}

func BlockChain() *blockChain {
	if b == nil {
		once.Do(func() {
			b = &blockChain{
				Height: 0,
			}
			persistBlockChain := db.BlockChain()
			if persistBlockChain == nil {
				b.AddBlock()
			} else {
				b.restore(persistBlockChain)
			}
		})
	}
	return b
}
