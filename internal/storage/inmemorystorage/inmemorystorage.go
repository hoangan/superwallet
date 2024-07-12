package inmemorystorage

import (
	"encoding/json"
	"fmt"
	"math/big"

	m "github.com/hoangan/superwallet/internal/models"
	inmemorydb "github.com/hoangan/superwallet/internal/storage/inmemorystorage/inmemorydatabase"
)

const (
	SubscribeAddressed = "subscribed_addresses"
	IndexedBlockNumber = "indexed_block_number"
)

type InMemoryStorage struct {
	db *inmemorydb.InMemoryDatabase
}

func New() (*InMemoryStorage, error) {
	storage := &InMemoryStorage{
		db: inmemorydb.New(),
	}

	// Initialize the database with subscribed addresses and their balances storage.
	// also easy to get list of subscribed addresses and their cache balances.
	addresses := make(map[string]big.Int)
	addressesBytes, err := json.Marshal(addresses)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize the database: %v", err)
	}

	if err = storage.db.Set(SubscribeAddressed, addressesBytes); err != nil {
		return nil, fmt.Errorf("failed to initialize the database: %v", err)
	}

	return storage, nil
}

func (s *InMemoryStorage) GetAddressesWithBalances() (map[string]big.Int, error) {
	var addresses map[string]big.Int

	addressesBytes, err := s.db.Get(SubscribeAddressed)
	if err != nil {
		return nil, fmt.Errorf("failed to get addresses: %v", err)
	}

	if err := json.Unmarshal(addressesBytes, &addresses); err != nil {
		return nil, fmt.Errorf("failed to get addresses: %v", err)
	}

	return addresses, nil
}

// SubscribeAddress adds a new address to the list of subscribed addresses.
// Addresses should be stored in address table in the case of sql database,
// with coin_id, chain_id, ticker, etc.
// For simplicity in the case of in-memory storage, we assume all addresses are for ETH native coin.
// Could be extended by adding the coind_id in the back of the address, e.g: address:coin_id.
func (s *InMemoryStorage) SubscribeAddress(address string) error {
	if _, err := s.db.Get(address); err == nil {
		return nil
	}

	// Get the current subscribed addresses and add the new address to the list.
	if addresses, err := s.GetAddressesWithBalances(); err != nil {
		return fmt.Errorf("failed to subscribe addresses: %v", err)
	} else {
		if _, ok := addresses[address]; ok {
			return nil
		}

		addresses[address] = *big.NewInt(0)
		addressesBytes, err := json.Marshal(addresses)
		if err != nil {
			return fmt.Errorf("failed to subscribe address: %v", err)
		}

		if err := s.db.Set(SubscribeAddressed, addressesBytes); err != nil {
			return fmt.Errorf("failed to subscribe address: %v", err)
		}
	}

	// Save the address and its future transactions's hashes in the database.
	// New subscribed address has no transactions yet.
	err := s.db.Set(address, []byte{})
	return err
}

func (s *InMemoryStorage) AddAddressTransaction(address string, txn m.Transaction) error {
	// Get the list of tx hash of the subscribed address.
	addressTxsBytes, err := s.db.Get(address)
	if err != nil {
		return fmt.Errorf("subscribed address does not exist: %v", err)
	}

	var addressTxs []string
	if err := json.Unmarshal(addressTxsBytes, &addressTxs); err != nil {
		return fmt.Errorf("failed to add address transaction: %v", err)
	}

	// Add the new txn hash to the list.
	addressTxs = append(addressTxs, txn.Hash)

	// Store the updated list back to the database.
	if addressTxsBytes, err := json.Marshal(addressTxs); err != nil {
		return fmt.Errorf("failed to add address transaction: %v", err)
	} else if err := s.db.Set(address, addressTxsBytes); err != nil {
		return fmt.Errorf("failed to add address transaction: %v", err)
	}

	// Store the txn only once, multiple addresses can have the same txn.
	// It's common for exchange to batch their withdrawals into a single transaction.
	if _, err := s.db.Get(txn.Hash); err == nil {
		return nil
	} else if err != inmemorydb.ErrNotFound {
		return fmt.Errorf("failed to add address transaction: %v", err)
	}

	if addressTxBytes, err := json.Marshal(addressTxs); err != nil {
		return fmt.Errorf("failed to add address transaction: %v", err)
	} else {
		if err := s.db.Set(txn.Hash, addressTxBytes); err != nil {
			return fmt.Errorf("failed to add address transaction: %v", err)
		}
	}

	return nil
}

func (s *InMemoryStorage) GetTransactionsByAddress(address string) ([]m.Transaction, error) {
	// Get the list of tx hash of the subscribed address.
	addressTxsBytes, err := s.db.Get(address)
	if err != nil {
		return nil, fmt.Errorf("subscribed address does not exist: %v", err)
	}

	var addressTxs []string
	if err := json.Unmarshal(addressTxsBytes, &addressTxs); err != nil {
		return nil, fmt.Errorf("failed to get address transactions: %v", err)
	}

	// Get the transactions by their hashes.
	var txns []m.Transaction
	for _, hash := range addressTxs {
		txBytes, err := s.db.Get(hash)
		if err != nil {
			return nil, fmt.Errorf("failed to get address transactions: %v", err)
		}

		var txn m.Transaction
		if err := json.Unmarshal(txBytes, &txn); err != nil {
			return nil, fmt.Errorf("failed to get address transactions: %v", err)
		}

		txns = append(txns, txn)
	}

	return txns, nil
}
