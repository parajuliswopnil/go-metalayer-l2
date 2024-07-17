package database

import (
	"io/fs"
	"math/big"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestNewDatabase(t *testing.T) {
	database := NewDatabase(dbPath)
	defer cleanup()
	assert.Equal(t, database.(*Database).Path, dbPath)
	assert.NotNil(t, database.(*Database).db)
}

func TestInitializeDatabase(t *testing.T) {
	var err error
	db, err = bolt.Open(dbPath, fs.FileMode(0600), nil)
	defer cleanup()
	assert.NoError(t, err)
	database := &Database{
		db: db,
	}

	err = database.InitializeDatabase()
	assert.NoError(t, err)

	database.db.View(func(tx *bolt.Tx) error {
		balanceBucket := tx.Bucket([]byte(accountBalanceBucket))
		assert.NotNil(t, balanceBucket)
		globalStateRootBucket := tx.Bucket([]byte(globalStateRootBucket))
		assert.NotNil(t, globalStateRootBucket)
		return nil
	})
}

func TestStoreAccountBalance(t *testing.T) {
	var err error
	db, err = bolt.Open(dbPath, fs.FileMode(0600), nil)
	defer cleanup()
	assert.NoError(t, err)
	database := &Database{
		db: db,
	}

	err = database.InitializeDatabase()
	assert.NoError(t, err)


	err = database.StoreAccountBalance("ethereum", "eth", common.HexToAddress("0xA1a98D3CFED036c9Ea921cCb046984EE36A8B849"), big.NewInt(100))
	assert.NoError(t, err)
}

func TestGetAccountBalance(t *testing.T) {
	var err error
	db, err = bolt.Open(dbPath, fs.FileMode(0600), nil)
	defer cleanup()
	assert.NoError(t, err)
	database := &Database{
		db: db,
	}

	err = database.InitializeDatabase()
	assert.NoError(t, err)

	err = database.StoreAccountBalance("ethereum", "eth", common.HexToAddress("0xA1a98D3CFED036c9Ea921cCb046984EE36A8B849"), big.NewInt(100))
	assert.NoError(t, err)

	balance := database.GetAccountBalance("ethereum", "eth", common.HexToAddress("0xA1a98D3CFED036c9Ea921cCb046984EE36A8B849"))
	cmp := balance.Cmp(big.NewInt(100))
	assert.Zero(t, cmp)
}
