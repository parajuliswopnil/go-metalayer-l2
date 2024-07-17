package database

import (
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func makeBalanceBucketKey(chain, token string, address ethCommon.Address) []byte {
	return crypto.Keccak256([]byte(chain), []byte(token), address.Bytes())
}

func makeMerkleLeaves(key, value []byte) []byte {
	return crypto.Keccak256(key, value)
}