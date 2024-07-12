package eth_test

import (
	"context"
	"math/big"
	"reflect"
	"testing"

	"github.com/hoangan/superwallet/internal/eth"
	"github.com/hoangan/superwallet/internal/storage/inmemorystorage"
	"github.com/hoangan/superwallet/internal/testdata"
)

const (
	ethEndpoint     = "https://cloudflare-eth.com"
	fromBlockNumber = 20290107
)

func TestEthIndexer(t *testing.T) {
	storage, _ := inmemorystorage.New()
	ethIndexer, _ := eth.NewIndexer(context.Background(), ethEndpoint, storage, big.NewInt(fromBlockNumber))

	t.Run("Parse Transaction", func(t *testing.T) {
		txn, err := ethIndexer.ParseTransaction(testdata.RawTransaction1)
		if err != nil {
			t.Errorf("failed to parse transaction: %w", err)
		}

		if !reflect.DeepEqual(txn, testdata.Transaction1) {
			t.Errorf("failed to parse transaction")
		}
	})
}
