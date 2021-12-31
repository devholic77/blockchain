package db

import (
	"github.com/boltdb/bolt"
	"github.com/devholic77/duckcoin/utils"
)

const (
	dbName       = "blockchain.db"
	dataBucket   = "data"
	blocksBucket = "blocks"
	checkpoint   = "checkpoint"
)

var db *bolt.DB

func DB() *bolt.DB {
	if db == nil {
		// init
		dbPointer, err := bolt.Open(dbName, 0600, nil)
		utils.HandleErr(err)
		db = dbPointer
		err = db.Update(func(t *bolt.Tx) error {
			_, err = t.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleErr(err)

			_, err = t.CreateBucketIfNotExists([]byte(blocksBucket))
			utils.HandleErr(err)
			return nil
		})
		utils.HandleErr(err)
	}
	return db
}

func SaveBlock(hash string, data []byte) error {

	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		return bucket.Put([]byte(hash), data)
	})
	return err
}

func Block(hash string) []byte {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(blocksBucket))
		data = b.Get([]byte(hash))
		return nil
	})
	return data
}

func SaveBlockChain(data []byte) error {
	err := DB().Update(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		return bucket.Put([]byte(checkpoint), data)
	})
	return err
}

func BlockChain() []byte {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		b := t.Bucket([]byte(dataBucket))
		data = b.Get([]byte(checkpoint))
		return nil
	})
	return data
}
