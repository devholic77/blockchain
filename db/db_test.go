package db

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

func TestMain(m *testing.M) {
	exit := m.Run()
	os.Exit(exit)
	DB()
}

func TestDBExist(t *testing.T) {
	DB()
	assert.Equal(t, true, true)
}

type TestStruct struct {
	Data string
}

// func TestDBPut(t *testing.T) {
// 	fmt.Println("test Put")
// 	testStruct := TestStruct{"test"}
// 	buf, _ := json.Marshal(testStruct)
// 	Put([]byte("test"), buf)
// }
// func TestDBGet(t *testing.T) {
// 	fmt.Println("test Get")
// 	Get([]byte("test"))
// }
