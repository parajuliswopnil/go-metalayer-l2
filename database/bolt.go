package database

import (
	"encoding/binary"
	"io/fs"
	"math/big"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

func newDatabase(path string) {
	var err error
	db, err = bolt.Open(path, fs.FileMode(0600), nil)
	if err != nil {
		panic(err)
	}
}

func newEntry(nameSpace, key, value []byte) error {
	var dbError error
	db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(nameSpace)
		if err != nil {
			dbError = err
		}
		err = bucket.Put(key, value)
		if err != nil {
			dbError = err
		}
		return nil
	})
	return dbError
}

func getEntry(nameSpace, key []byte) ([]byte, error) {
	var value []byte
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(nameSpace)
		value = bucket.Get(key)
		return nil
	})
	return value, nil
}

type KeyConstraints interface {
	uint64 | *big.Int | ~string | ~[]byte
}

func convertToByte[K KeyConstraints](input K) []byte {
	var i interface{} = input
	switch v := i.(type) {
	case uint64:
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, v)
		return b
	case *big.Int:
		return v.Bytes()
	case string:
		return []byte(v)
	case []byte:
		return v
	default:
		return nil
	}
}

func getAllEntries(namespace []byte) [][]byte {
	var entries [][]byte
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(namespace)
		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			entries = append(entries, v)
		}
		return nil
	})
	return entries
}

func closeDB() error {
	return db.Close()
}
