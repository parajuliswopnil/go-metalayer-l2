package database

import (
	"fmt"
	"io/fs"
	"math/big"
	"strconv"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/cedro-finance/metalayer-sequencer/merkle"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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

func TestStorageMerkleRoot(t *testing.T) {
	var err error
	db, err = bolt.Open(dbPath, fs.FileMode(0600), nil)
	defer cleanup()
	assert.NoError(t, err)
	database := &Database{
		db: db,
	}

	err = database.InitializeDatabase()
	assert.NoError(t, err)

	for i := range 8 {
		err = database.StoreAccountBalance("ethereum", "eth", common.HexToAddress("0xA1a98D3CFED036c9Ea921cCb046984EE36A8B84" + strconv.Itoa(i)), big.NewInt(100))
		assert.NoError(t, err)
	}

	root, err := database.GetGlobalStateRoot()
	assert.NoError(t, err)
	fmt.Println(root)
}

func TestVerifyMerkleProof(t *testing.T) {
	var err error
	db, err = bolt.Open(dbPath, fs.FileMode(0600), nil)
	defer cleanup()
	assert.NoError(t, err)
	database := &Database{
		db: db,
	}

	err = database.InitializeDatabase()
	assert.NoError(t, err)

	for i := range 8 {
		err = database.StoreAccountBalance("ethereum", "eth", common.HexToAddress("0xA1a98D3CFED036c9Ea921cCb046984EE36A8B84" + strconv.Itoa(i)), big.NewInt(100))
		assert.NoError(t, err)
	}

	root, err := database.GetGlobalStateRoot()
	assert.NoError(t, err)

	account := makeBalanceBucketKey("ethereum", "eth", common.HexToAddress("0xA1a98D3CFED036c9Ea921cCb046984EE36A8B840"))
	accLf := accountLeaf[hexutil.Encode(account)]
	hashedAccountLeaf := merkle.Hasher(accLf)
	proof := accountProof[hexutil.Encode(hashedAccountLeaf)]

	assert.True(t, database.VerifyMerkleProofForAccountBalance(hashedAccountLeaf, root, proof))
}
