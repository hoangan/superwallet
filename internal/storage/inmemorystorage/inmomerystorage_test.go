package inmemorystorage_test

import (
	"math/big"
	"testing"

	m "github.com/hoangan/superwallet/internal/models"
	"github.com/hoangan/superwallet/internal/storage/inmemorystorage"
)

func TestInMemoryStorage(t *testing.T) {
	storage, _ := inmemorystorage.New()
	address := "0x29182006a4967e9a50c0a66076da514993d3b4d4"
	storage.SubscribeAddress(address)
	if !storage.IsSubscribedAddress(address) {
		t.Errorf("failed to subscribe address")
	}

	err := storage.AddAddressTransaction(address, &m.Transaction{
		Hash:        "0x6523b98c957773ece31e36dfe4309df7cd6cd697f70b7f4df6d39fc008cc693a",
		Type:        big.NewInt(2),
		BlockHash:   "0xf20326ecb02332687c918de6df6c8b354ccdf8406ea1b276a4da07e22b072715",
		BlockNumber: big.NewInt(20290107),
		ChainId:     big.NewInt(1),
		Nonce:       big.NewInt(0),
		Gas:         big.NewInt(21000),
		GasPrice:    big.NewInt(1000000000),
		Transfers: []*m.Transfer{
			{
				CoinID: 1,
				From:   "0x4838B106FCe9647Bdf1E7877BF73cE8B0BAD5f97",
				To:     "0x29182006a4967e9a50c0a66076da514993d3b4d4",
			},
		},
	})

	if err != nil {
		t.Errorf("failed to add address transaction: %v", err)
	}

	transactions, err := storage.GetTransactionsByAddress(address)
	if err != nil {
		t.Errorf("failed to get transactions by address: %v", err)
	}

	if len(transactions) != 1 {
		t.Errorf("failed to get transactions by address")
	}
}
