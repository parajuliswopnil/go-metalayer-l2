package database

import (
	"encoding/binary"
	"fmt"
	"io/fs"
	"math/big"
	"os"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)

const dbPath = "bolt.db"

func cleanup() {
	db.Close()
	os.Remove(dbPath)
}
func TestNewDb(t *testing.T) {
	assert.Nil(t, db)
	newDatabase("bolt.db")
	defer cleanup()
	assert.NotNil(t, db)
}

func TestNewEntry(t *testing.T) {
	var err error
	db, err = bolt.Open(dbPath, fs.FileMode(0600), nil)
	defer cleanup()
	assert.NoError(t, err)
	nameSpace := []byte("newNamespace")
	key := []byte("key")
	value := []byte("value")

	err = newEntry(nameSpace, key, value)
	assert.NoError(t, err)
}

func TestGetEntry(t *testing.T) {
	var err error
	db, err = bolt.Open(dbPath, fs.FileMode(0600), nil)
	defer cleanup()
	assert.NoError(t, err)
	nameSpace := []byte("newNamespace")
	key := []byte("key")
	value := []byte("value")

	err = newEntry(nameSpace, key, value)
	assert.NoError(t, err)

	dbValue, err := getEntry(nameSpace, key)
	assert.NoError(t, err)
	assert.Equal(t, value, dbValue)
}

func TestConvertToByte(t *testing.T) {
	uintValue := uint64(123)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uintValue)
	assert.Equal(t, b, convertToByte(uintValue))

	bigIntValue := big.NewInt(123)
	bigIntBytes := bigIntValue.Bytes()
	assert.Equal(t, bigIntBytes, convertToByte(bigIntValue))

	stringValue := "123"
	stringByte := []byte(stringValue)
	assert.Equal(t, stringByte, convertToByte(stringValue))

	byteValue := []byte{1, 2, 3}
	assert.Equal(t, byteValue, convertToByte(byteValue))
}

func TestGetAllEntries(t *testing.T) {
	var err error
	db, err = bolt.Open(dbPath, fs.FileMode(0600), nil)
	defer cleanup()
	assert.NoError(t, err)
	nameSpace := []byte("newNamespace")
	key := []byte("key1")
	value := []byte("value1")

	key1 := []byte("key")
	value1 := []byte("value")

	err = newEntry(nameSpace, key, value)
	assert.NoError(t, err)

	err = newEntry(nameSpace, key1, value1)
	assert.NoError(t, err)

	dbValue := getAllEntries(nameSpace)

	for _, v := range dbValue {
		fmt.Println(string(v))
	}
}

func TestCloseDB(t *testing.T) {
	var err error
	db, err = bolt.Open(dbPath, fs.FileMode(0600), nil)
	defer os.Remove(dbPath)
	assert.NoError(t, err)
	path := db.Path()
	assert.Equal(t, path, dbPath)
	closeDB()
	path = db.Path()
	assert.Equal(t, path, "")
}
