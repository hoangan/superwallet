package internal

import (
	"math/big"

	m "github.com/hoangan/superwallet/internal/models"
)

// Indexer is the interface that wraps the basic methods for transaction parser
// Couble be extended with more services like: hot wallet withdraw, fund sweep in the case custodial wallet
type Indexer interface {
	// start indexer
	Start()

	// last parsed block number
	GetCurrentBlock() *big.Int

	// add address to observer
	SubscribeAddress(address string) error

	// list of inbound or outbound transactions for an address
	GetTransactions(address string) ([]*m.Transaction, error)
}
