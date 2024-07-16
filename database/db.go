package database

import (
	"math/big"

	"github.com/boltdb/bolt"
	ethCommon "github.com/ethereum/go-ethereum/common"
)

const (
	accountBalanceBucket = "account_balance_bucket"
	globalStateRootBucket = "global_state_root_bucket"
)

type IDatabase interface {
	StoreAccountBalance()
	GetAccountBalance(chain, token string, address ethCommon.Address) *big.Int
}

type Database struct {
	Path string 
	db *bolt.DB
}

func NewDatabase(path string) *Database {
	newDatabase(path)
	return &Database{
		Path: path,
		db: db,
	}
}

// creates buckets for storing
// 1. Account Balances
// 2. Previous Global State root 
// 3. Current Global State rooot
func (d *Database) InitializeDatabase() error {
	var err error 
	d.db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(accountBalanceBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(globalStateRootBucket))
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// store account information 
// key: Hash(chain + token + address)
// value: *big.Int
func (d *Database) StoreAccountBalance() {}

// get balance information of account account 
func (d *Database) GetAccountBalance(chain, token string, address ethCommon.Address) *big.Int {return ethCommon.Big0}

