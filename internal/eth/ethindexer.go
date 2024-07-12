package eth

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/hoangan/superwallet/internal/eth/rpc"
	"github.com/hoangan/superwallet/internal/storage"
)

const (
	defaultFromBlockNumber = 15537393
	blockTime              = 15 // seconds
)

type EthIndexer struct {
	ctx                 context.Context
	ticker              *time.Ticker
	client              *rpc.EthClient
	currentIndexedBlock *big.Int
	storage             *storage.Storage
	once                sync.Once
	wg                  sync.WaitGroup
}

func New(ctx context.Context, client *rpc.EthClient, storage *storage.Storage, fromBlockNumber *big.Int) (*EthIndexer, error) {
	var err error
	var currentIndexedBlock *big.Int
	if fromBlockNumber == nil {
		if currentIndexedBlock, err = (*storage).GetIndexedBlockNumber(); err != nil {
			currentIndexedBlock = big.NewInt(defaultFromBlockNumber)
		}

	}

	return &EthIndexer{
		ctx:                 ctx,
		ticker:              time.NewTicker(blockTime * time.Second),
		client:              client,
		currentIndexedBlock: currentIndexedBlock,
		storage:             storage,
	}, nil
}

func (i *EthIndexer) StartIndexing(ctx context.Context) error {
	i.once.Do(func() {
		i.wg.Add(1)
		go func() {
			defer i.wg.Done()
			for {
				select {
				case <-ctx.Done():
					i.Stop()
					return
				default:
					i.ticker.Stop()
					// get the latest block number
					// latestBlock, err := i.client.GetLatestBlockNumber()
					// if err != nil {
					// 	fmt.Printf("failed to get the latest block number: %v", err)
					// 	time.Sleep(blockTime * time.Second)
					// 	i.ticker.Reset(blockTime * time.Second)
					// }

					// if i.currentIndexedBlock.Cmp(latestBlock) >= 0 {
					// 	continue
					// }

					// // fetch the block by number
					// for i.currentIndexedBlock.Cmp(latestBlock) < 0 {
					// 	next
					// }

					i.ticker.Reset(blockTime * time.Second)

				}
			}
		}()

	})
	return nil
}

func (i *EthIndexer) Stop() {
	i.ticker.Stop()
	i.wg.Wait()
}
