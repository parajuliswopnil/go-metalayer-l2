package database

import (
	"io/fs"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)

func TestNewDatabase(t *testing.T) {
	database := NewDatabase(dbPath)
	defer cleanup()
	assert.Equal(t, database.Path, dbPath)
	assert.NotNil(t, database.db)
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
