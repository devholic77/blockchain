package utils

import (
	"bytes"
	"encoding/json"
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
	decoder.Decode(&v)
}
