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
	if err := storage.encodeAndSave(SubscribeAddressed, make(map[string]big.Int)); err != nil {
		return nil, fmt.Errorf("failed to initialize the database: %w", err)
	}

	// Initialize the database with the indexed block number.
	// This is used to keep track of the last indexed block number.
	if err := storage.encodeAndSave(IndexedBlockNumber, big.NewInt(0)); err != nil {
		return nil, fmt.Errorf("failed to initialize the database: %w", err)
	}

	return storage, nil
}

func (s *InMemoryStorage) GetIndexedBlockNumber() (*big.Int, error) {
	indexedBlockNumberBytes, err := s.db.Get(IndexedBlockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get indexed block number from db: %w", err)
	}

	var indexedBlockNumber big.Int
	if err := json.Unmarshal(indexedBlockNumberBytes, &indexedBlockNumber); err != nil {
		return nil, fmt.Errorf("failed to get indexed block number from db: %w", err)
	}

	return &indexedBlockNumber, nil
}

func (s *InMemoryStorage) SaveIndexedBlockNumber(indexedBlockNumber *big.Int) error {
	if err := s.encodeAndSave(IndexedBlockNumber, indexedBlockNumber); err != nil {
		return fmt.Errorf("failed to save indexed block number: %w", err)
	}

	return nil
}

func (s *InMemoryStorage) GetAddressesWithBalances() (map[string]*big.Int, error) {
	var addresses map[string]*big.Int

	addressesBytes, err := s.db.Get(SubscribeAddressed)
	if err != nil {
		return nil, fmt.Errorf("failed to get addresses: %w", err)
	}

	if err := json.Unmarshal(addressesBytes, &addresses); err != nil {
		return nil, fmt.Errorf("failed to get addresses: %w", err)
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
		return fmt.Errorf("failed to subscribe addresses: %w", err)
	} else {
		if _, ok := addresses[address]; ok {
			return nil
		}

		// Add the new address to the collection of subscribed addresses.
		// with initial balance of 0.
		addresses[address] = big.NewInt(0)
		if err := s.encodeAndSave(SubscribeAddressed, addresses); err != nil {
			return fmt.Errorf("failed to subscribe address: %w", err)
		}
	}

	// Save the address and its future transactions's hash list in the database.
	// New subscribed address has no transactions yet.
	if err := s.encodeAndSave(address, []string{}); err != nil {
		return fmt.Errorf("failed to subscribe address: %w", err)
	}

	return nil
}

func (s *InMemoryStorage) AddAddressTransaction(address string, txn *m.Transaction) error {
	// Store the txn only once, multiple addresses can have the same txn.
	// It's common for exchange to batch their withdrawals into a single transaction.
	if _, err := s.db.Get(txn.Hash); err == nil {
		return nil
	} else if err != inmemorydb.ErrNotFound { //other error, e.g.: db closed
		return fmt.Errorf("failed to add address transaction: %w", err)
	}

	if err := s.encodeAndSave(txn.Hash, txn); err != nil {
		return fmt.Errorf("failed to save transaction: %w", err)
	}

	// Get the list of tx hash of the subscribed address.
	addressTxHashesBytes, err := s.db.Get(address)
	if err != nil {
		return fmt.Errorf("subscribed address does not exist: %w", err)
	}

	var addressTxHashes []string
	if err := json.Unmarshal(addressTxHashesBytes, &addressTxHashes); err != nil {
		return fmt.Errorf("failed to fetch current address tx hash list: %w", err)
	}

	// Check if the txn hash already exists in the list.
	for _, hash := range addressTxHashes {
		if hash == txn.Hash {
			return nil
		}
	}

	// Add the new txn hash to the list.
	addressTxHashes = append(addressTxHashes, txn.Hash)
	if err := s.encodeAndSave(address, addressTxHashes); err != nil {
		return fmt.Errorf("failed to add address transaction: %w", err)
	}

	return nil
}

func (s *InMemoryStorage) GetTransactionsByAddress(address string) ([]*m.Transaction, error) {
	// Get the list of tx hash of the subscribed address.
	addressTxsBytes, err := s.db.Get(address)
	if err != nil {
		return nil, fmt.Errorf("subscribed address does not exist: %w", err)
	}

	var addressTxs []string
	if err := json.Unmarshal(addressTxsBytes, &addressTxs); err != nil {
		return nil, fmt.Errorf("failed to get address tx hash list: %w", err)
	}

	// Get the transactions by their hashes.
	var txns []*m.Transaction
	for _, hash := range addressTxs {
		txBytes, err := s.db.Get(hash)
		if err != nil {
			return nil, fmt.Errorf("failed to get transaction by hash: %w", err)
		}

		var txn m.Transaction
		if err := json.Unmarshal(txBytes, &txn); err != nil {
			fmt.Printf("failed to load address transaction: %s\n", string(txBytes))
			return nil, fmt.Errorf("failed to load address transaction: %w", err)
		}

		txns = append(txns, &txn)
	}

	return txns, nil
}

func (s *InMemoryStorage) IsSubscribedAddress(address string) bool {
	if _, err := s.db.Get(address); err != nil {
		return false
	}

	return true
}

// encodeAndSave marshal any value data type and saves it to the database as bytes.
func (s *InMemoryStorage) encodeAndSave(key string, value interface{}) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to save value: %w", err)
	}

	if err := s.db.Set(key, valueBytes); err != nil {
		return fmt.Errorf("failed to save value: %w", err)
	}

	return nil
}
