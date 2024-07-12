package storage

import (
	"math/big"

	m "github.com/hoangan/superwallet/internal/models"
)

// Storage interface is the interface that wraps the basic methods for a storage.
type Storage interface {
	SubscribeAddress(address string) error
	GetTransactionsByAddress(address string) ([]*m.Transaction, error)
	AddAddressTransaction(address string, tx *m.Transaction) error
	GetAddressesWithBalances() (map[string]*big.Int, error)
	SaveIndexedBlockNumber(indexedBlockNumber *big.Int) error
	GetIndexedBlockNumber() (*big.Int, error)
	IsSubscribedAddress(address string) bool
}
