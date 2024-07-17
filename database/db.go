package database

import (
	"math/big"

	"github.com/boltdb/bolt"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

const (
	accountBalanceBucket  = "account_balance_bucket"
	globalStateRootBucket = "global_state_root_bucket"
)

type IDatabase interface {
	InitializeDatabase() error
	StoreAccountBalance(string, string, ethCommon.Address, *big.Int) error
	GetAccountBalance(string, string, ethCommon.Address) *big.Int
}

type Database struct {
	Path string
	db   *bolt.DB
}

func NewDatabase(path string) IDatabase {
	newDatabase(path)
	return &Database{
		Path: path,
		db:   db,
	}
}

// creates buckets for storing
// 1. Account Balances
// 2. Previous Global State root
// 3. Current Global State rooot
func (d *Database) InitializeDatabase() error {
	return d.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(accountBalanceBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(globalStateRootBucket))
		if err != nil {
			return err
		}
		return nil
	})
}

// store account information
// key: Hash(chain + token + address)
// value: *big.Int
func (d *Database) StoreAccountBalance(chain, token string, address ethCommon.Address, value *big.Int) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		balanceBucket := tx.Bucket([]byte(accountBalanceBucket))
		if balanceBucket == nil {
			return bolt.ErrBucketNotFound
		}
		key := makeBalanceBucketKey(chain, token, address)
		err := balanceBucket.Put(key, value.Bytes())
		if err != nil {
			return err
		}
		return nil
	})
}

// get balance information of account account
func (d *Database) GetAccountBalance(chain, token string, address ethCommon.Address) *big.Int {
	balance := new(big.Int)
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(accountBalanceBucket))
		if bucket == nil {
			return bolt.ErrBucketNotFound
		}
		key := makeBalanceBucketKey(chain, token, address)

		value := bucket.Get(key)
		if value == nil {
			balance = ethCommon.Big0
			return nil
		}
		balance = balance.SetBytes(value)
		return nil
	})
	return balance
}

func (d *Database) getBalanceStorageLeaves() ([][]byte, error) {
	var leaves [][]byte
	_ = leaves
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(accountBalanceBucket))
		if bucket == nil {
			return bolt.ErrBucketNotFound
		}

		cursor := bucket.Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			leaf := makeMerkleLeaves(k, v)
			leaves = append(leaves, leaf)
		}
		return nil
	})
	return leaves, err
}
