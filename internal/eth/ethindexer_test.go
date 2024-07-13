package eth_test

import (
	"context"
	"math/big"
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
			t.Errorf("failed to parse transaction: %v", err)
		}

		if txn.Hash != testdata.Transaction1.Hash {
			t.Errorf("failed to parse transaction hash: %s", txn.Hash)
		}

		if txn.Type.Cmp(testdata.Transaction1.Type) != 0 {
			t.Errorf("failed to parse transaction type: %s", txn.Type)
		}

		if txn.BlockHash != testdata.Transaction1.BlockHash {
			t.Errorf("failed to parse transaction block hash: %s", txn.BlockHash)
		}

		if txn.BlockNumber.Cmp(testdata.Transaction1.BlockNumber) != 0 {
			t.Errorf("failed to parse transaction block number: %s", txn.BlockNumber)
		}

		if txn.ChainId.Cmp(testdata.Transaction1.ChainId) != 0 {
			t.Errorf("failed to parse transaction chain id: %s", txn.ChainId)
		}

		if txn.Nonce.Cmp(testdata.Transaction1.Nonce) != 0 {
			t.Errorf("failed to parse transaction nonce: %s", txn.Nonce)
		}

		if txn.Gas.Cmp(testdata.Transaction1.Gas) != 0 {
			t.Errorf("failed to parse transaction gas: %s", txn.Gas)
		}

		if txn.GasPrice.Cmp(testdata.Transaction1.GasPrice) != 0 {
			t.Errorf("failed to parse transaction gas price: %s", txn.GasPrice)
		}

		if txn.From != testdata.Transaction1.From {
			t.Errorf("failed to parse transaction from: %s", txn.From)
		}

		if txn.To != testdata.Transaction1.To {
			t.Errorf("failed to parse transaction to: %s", txn.To)
		}

		if txn.Value.Cmp(testdata.Transaction1.Value) != 0 {
			t.Errorf("failed to parse transaction value: %s", txn.Value)
		}

		if len(txn.Transfers) != len(testdata.Transaction1.Transfers) {
			t.Errorf("failed to parse transaction transfers: %v", txn.Transfers)
		}

		for i, transfer := range txn.Transfers {
			if transfer.CoinID != testdata.Transaction1.Transfers[i].CoinID {
				t.Errorf("failed to parse transfer coin id: %d", transfer.CoinID)
			}

			if transfer.Ticker != testdata.Transaction1.Transfers[i].Ticker {
				t.Errorf("failed to parse transfer ticker: %s", transfer.Ticker)
			}

			if transfer.From != testdata.Transaction1.Transfers[i].From {
				t.Errorf("failed to parse transfer from: %s", transfer.From)
			}

			if transfer.To != testdata.Transaction1.Transfers[i].To {
				t.Errorf("failed to parse transfer to: %s", transfer.To)
			}

			if transfer.Value.Cmp(testdata.Transaction1.Transfers[i].Value) != 0 {
				t.Errorf("failed to parse transfer value: %s", transfer.Value)
			}
		}
	})
}
