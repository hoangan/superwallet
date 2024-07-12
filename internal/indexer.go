package internal

import (
	m "github.com/hoangan/superwallet/internal/models"
)

// Indexer is the interface that wraps the basic methods for transaction parser
// Couble be extended with more services like: hot wallet withdraw, fund sweep in the case custodial wallet
type Indexer interface {
	// start indexer
	Start() error

	// last parsed block number
	GetCurrentIndexedBlock() (int64, error)

	// add address to observer
	SubscribeAddress(address string) error

	// list of inbound or outbound transactions for an address
	GetTransactions(address string) ([]m.Transaction, error)
}
