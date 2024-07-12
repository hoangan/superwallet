package eth

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/hoangan/superwallet/internal/eth/rpc"
	m "github.com/hoangan/superwallet/internal/models"
	"github.com/hoangan/superwallet/internal/storage"
	"github.com/hoangan/superwallet/pkg/enccode/hexencoder"
)

const (
	DefaultFromBlockNumber = 15537393
	blockTime              = 15    // seconds
	chainId                = 1     // mainnet
	coinId                 = 1     // ethereum
	coinTicker             = "ETH" // ethereum
)

type EthIndexer struct {
	ctx                 context.Context
	ticker              *time.Ticker
	client              *rpc.EthClient
	currentIndexedBlock *big.Int
	storage             storage.Storage
	once                sync.Once
	wg                  sync.WaitGroup
}

func NewIndexer(ctx context.Context, endpoint string, storage storage.Storage, fromBlockNumber *big.Int) (*EthIndexer, error) {
	var err error
	var currentIndexedBlock *big.Int
	if fromBlockNumber == nil {
		// in production system, load the last indexed block from the database
		// for case where by the system is restarted
		if currentIndexedBlock, err = storage.GetIndexedBlockNumber(); err != nil {
			currentIndexedBlock = big.NewInt(DefaultFromBlockNumber)
		}
	} else {
		currentIndexedBlock = fromBlockNumber.Sub(fromBlockNumber, big.NewInt(1))
	}

	return &EthIndexer{
		ctx:                 ctx,
		ticker:              time.NewTicker(blockTime * time.Second),
		client:              rpc.NewEthClient(endpoint),
		currentIndexedBlock: currentIndexedBlock,
		storage:             storage,
	}, nil
}

func (i *EthIndexer) Start() {
	i.once.Do(func() {
		i.wg.Add(1)
		go func() {
			defer i.wg.Done()
			for {
				select {
				case <-i.ctx.Done():
					i.Stop()
					return
				default:
					i.ticker.Stop()
					latestRawBlock, err := i.client.GetLatestBlock()

					if err != nil {
						fmt.Printf("failed to get latest block: %v. Retry in %ds...\n", err, blockTime)

						// In case of error, node is not reachable, wait for a block time before retrying
						time.Sleep(blockTime * time.Second)
						continue
					}

					latestBlockNumber, err := hexencoder.HexToDecimal(latestRawBlock.Number)
					if err != nil {
						fmt.Printf("failed to parse latest block number %s: %v. Retry in %ds...\n", latestRawBlock.Number, err, blockTime)

						// In case of error, wait for a block time before retrying
						// node does return gibberish data when it is faulty sometime
						time.Sleep(blockTime * time.Second)
						continue
					}

					if i.currentIndexedBlock.Cmp(latestBlockNumber) < 0 {
						currentBlockNumber := i.currentIndexedBlock.Add(i.currentIndexedBlock, big.NewInt(1))

						currentRawBlock, err := i.client.GetBlockByNumber(currentBlockNumber)
						if err != nil {
							fmt.Printf("failed to get block by number: %v. Retry in %ds...\n", err, blockTime)

							// In case of error, wait for a block time before retrying
							time.Sleep(blockTime * time.Second)
							continue
						}

						for _, rawTx := range currentRawBlock.Transactions {
							tx, err := i.ParseTransaction(rawTx)
							if err != nil {
								fmt.Printf("failed to parse transaction: %v\n", err)
								continue
							}
							err = i.SaveSubscibedAddressTransaction(tx)
							if err != nil {
								fmt.Printf("failed to save subscribed address transaction: %v\n", err)
								continue
							}
						}

						i.currentIndexedBlock = currentBlockNumber
					}

					i.ticker.Reset(blockTime * time.Second)

				}
			}
		}()

	})
}

func (i *EthIndexer) ParseTransaction(rawTxn *rpc.RawTransaction) (*m.Transaction, error) {
	var err error
	tx := &m.Transaction{}
	transfers := []*m.Transfer{}

	// parse raw transaction to transaction for internal use
	tx.Hash = rawTxn.Hash

	if tx.Type, err = hexencoder.HexToDecimal(rawTxn.Type); err != nil {
		return nil, fmt.Errorf("failed to parse type: %v", err)
	}

	if tx.BlockNumber, err = hexencoder.HexToDecimal(rawTxn.BlockNumber); err != nil {
		return nil, fmt.Errorf("failed to parse block number: %v", err)
	}
	tx.BlockHash = rawTxn.BlockHash
	tx.From = rawTxn.From
	tx.To = rawTxn.To

	if tx.Value, err = hexencoder.HexToDecimal(rawTxn.Value); err != nil {
		return nil, fmt.Errorf("failed to parse value: %v", err)
	}

	// In the case of contract call, the value is 0
	transfers = append(transfers, &m.Transfer{
		CoinID: coinId,
		Ticker: coinTicker,
		From:   tx.From,
		To:     tx.To,
		Value:  tx.Value,
	})

	tx.Input = rawTxn.Input
	if len(tx.Input) > 2 {
		// TODO: parse internal transactions using trace_transaction
		// However, the node api provided by cloudflare does not support this method
	}

	if tx.ChainId, err = hexencoder.HexToDecimal(rawTxn.ChainId); err != nil {
		// default chain id to mainnet for now
		// legacy transactions do not have this field
		tx.ChainId = big.NewInt(chainId)
		// return nil, fmt.Errorf("failed to parse chain id: %v\n %+v", err, rawTxn)
	}

	if tx.Nonce, err = hexencoder.HexToDecimal(rawTxn.Nonce); err != nil {
		return nil, fmt.Errorf("failed to parse nonce: %v", err)
	}

	if tx.Gas, err = hexencoder.HexToDecimal(rawTxn.Gas); err != nil {
		return nil, fmt.Errorf("failed to parse gas: %v", err)
	}

	if tx.GasPrice, err = hexencoder.HexToDecimal(rawTxn.GasPrice); err != nil {
		return nil, fmt.Errorf("failed to parse gas price: %v", err)
	}

	tx.Transfers = transfers

	return tx, nil
}

// Check if the transaction contains subscribed address
// then save it to the database
func (i *EthIndexer) SaveSubscibedAddressTransaction(tx *m.Transaction) error {
	for _, transfer := range tx.Transfers {
		if i.storage.IsSubscribedAddress(transfer.From) {
			if err := i.storage.AddAddressTransaction(transfer.From, tx); err != nil {
				fmt.Printf("failed to save transaction subscribed address %s : %v", transfer.From, err)
			} else {
				fmt.Printf("saved transaction for subscribed address: %s hash: %s\n", transfer.From, tx.Hash)
			}
		}

		if i.storage.IsSubscribedAddress(transfer.To) {
			if err := i.storage.AddAddressTransaction(transfer.To, tx); err != nil {
				fmt.Printf("failed to save transaction subscribed address %s : %+v", transfer.To, err)
			} else {
				fmt.Printf("saved transaction for subscribed address: %s hash: %s\n", transfer.To, tx.Hash)
			}
		}
	}

	return nil
}

func (i *EthIndexer) GetCurrentBlock() *big.Int {
	return i.currentIndexedBlock
}

func (i *EthIndexer) GetTransactions(address string) ([]*m.Transaction, error) {
	return i.storage.GetTransactionsByAddress(address)
}

func (i *EthIndexer) Stop() {
	i.ticker.Stop()

	// wait for the indexer to finish processing the current block
	i.wg.Wait()
}
