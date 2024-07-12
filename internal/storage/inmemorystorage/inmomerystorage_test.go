package inmemorystorage_test

import (
	"testing"

	"github.com/hoangan/superwallet/internal/storage/inmemorystorage"
	"github.com/hoangan/superwallet/internal/testdata"
)

func TestInMemoryStorage(t *testing.T) {
	storage, _ := inmemorystorage.New()
	address := "0x29182006a4967e9a50c0a66076da514993d3b4d4"

	t.Run("Subscribe Address", func(t *testing.T) {
		err := storage.SubscribeAddress(address)
		if err != nil {
			t.Errorf("failed to subscribe address: %w", err)
		}

		if !storage.IsSubscribedAddress(address) {
			t.Errorf("failed to subscribe address")
		}
	})

	t.Run("Add Address Transaction", func(t *testing.T) {
		err := storage.AddAddressTransaction(address, testdata.Transaction1)
		if err != nil {
			t.Errorf("failed to add address transaction: %w", err)
		}

	})

	t.Run("Get Transactions By Address", func(t *testing.T) {
		transactions, err := storage.GetTransactionsByAddress(address)
		if err != nil {
			t.Errorf("failed to get transactions by address: %w", err)
		}

		if len(transactions) != 1 {
			t.Errorf("failed to get transactions by address txs count: %d", len(transactions))
		}
	})
}
