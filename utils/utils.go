package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
)

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToBytes(i interface{}) []byte {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	HandleErr(encoder.Encode(i))
	return buf.Bytes()
}

func FromBytes(data []byte, v interface{}) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.Decode(v)
}

func Hash(i interface{}) string {
	s := fmt.Sprintf("%v", i)
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}
