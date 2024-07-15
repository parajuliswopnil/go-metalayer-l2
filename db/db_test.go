package db

import (
	"encoding/binary"
	"io/fs"
	"math/big"
	"os"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)

const dbPath = "bolt.db"

func cleanup() {
	DB.Close()
	os.Remove(dbPath)
}
func TestNewDb(t *testing.T) {
	assert.Nil(t, DB)
	NewDatabase("bolt.db")
	defer cleanup()
	assert.NotNil(t, DB)
}

func TestNewEntry(t *testing.T) {
	var err error
	DB, err = bolt.Open(dbPath, fs.FileMode(0600), nil)
	defer cleanup()
	assert.NoError(t, err)
	nameSpace := []byte("newNamespace")
	key := []byte("key")
	value := []byte("value")

	err = NewEntry(nameSpace, key, value)
	assert.NoError(t, err)
}

func TestGetEntry(t *testing.T) {
	var err error
	DB, err = bolt.Open(dbPath, fs.FileMode(0600), nil)
	defer cleanup()
	assert.NoError(t, err)
	nameSpace := []byte("newNamespace")
	key := []byte("key")
	value := []byte("value")

	err = NewEntry(nameSpace, key, value)
	assert.NoError(t, err)

	dbValue, err := GetEntry(nameSpace, key)
	assert.NoError(t, err)
	assert.Equal(t, value, dbValue)
}

func TestConvertToByte(t *testing.T) {
	uintValue := uint64(123)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uintValue)
	assert.Equal(t, b, ConvertToByte(uintValue))

	bigIntValue := big.NewInt(123)
	bigIntBytes := bigIntValue.Bytes()
	assert.Equal(t, bigIntBytes, ConvertToByte(bigIntValue))

	stringValue := "123"
	stringByte := []byte(stringValue)
	assert.Equal(t, stringByte, ConvertToByte(stringValue))

	byteValue := []byte{1, 2, 3}
	assert.Equal(t, byteValue, ConvertToByte(byteValue))
}