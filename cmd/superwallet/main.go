package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/hoangan/superwallet/internal"
	"github.com/hoangan/superwallet/internal/eth"
	"github.com/hoangan/superwallet/internal/storage/inmemorystorage"
)

const (
	EthEndpoint = "https://cloudflare-eth.com"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("could not start the indexer: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var ethIndexer internal.Indexer

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fromBlockNumber := flag.Int("from-block", eth.DefaultFromBlockNumber, "from block number to start indexing")
	flag.Parse()

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)

	storage, err := inmemorystorage.New()
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}

	ethIndexer, err = eth.NewIndexer(ctx, EthEndpoint, storage, big.NewInt(int64(*fromBlockNumber)))
	if err != nil {
		return fmt.Errorf("failed to create eth indexer: %w", err)
	}

	ethIndexer.Start()

	fmt.Printf("Indexer started...\n")

	reader := bufio.NewReader(os.Stdin)

	{
	Quit:
		for {
			select {
			case <-terminate:
				fmt.Printf("Stopping indexer...\n")
				cancel()
				break Quit
			default:
				fmt.Printf("-> ")
				input, err := reader.ReadString('\n')
				if err != nil {
					fmt.Printf("failed to read input: %v\n", err)
				}
				input = strings.TrimSpace(input)
				args := strings.Split(input, " ")
				switch args[0] {
				case "\\q":
					fmt.Printf("Stopping indexer...\n")
					cancel()
					break Quit
				case "\\s":
					if len(args) < 2 {
						fmt.Printf("missing address\n")
						continue
					}
					address := args[1]
					if err := ethIndexer.SubscribeAddress(strings.ToLower(address)); err != nil {
						fmt.Printf("failed to subscribe address: %v\n", err)
					}
					fmt.Printf("address %s subscribed\n", address)
				case "\\a":
					if len(args) < 2 {
						fmt.Printf("missing address\n")
						continue
					}
					address := args[1]
					transactions, err := ethIndexer.GetTransactions(strings.ToLower(address))
					if err != nil {
						fmt.Printf("failed to get transactions: %v\n", err)
						continue
					}

					for _, tx := range transactions {
						txBytes, err := json.Marshal(tx)
						if err != nil {
							fmt.Printf("failed to marshal transaction: %+v\n", err)
							continue
						}
						fmt.Printf("%s\n\n", txBytes)
					}
				case "\\b":
					currentIndexedBlock := ethIndexer.GetCurrentBlock()
					fmt.Printf("current indexed block: %s\n", currentIndexedBlock.String())
				}
			}
		}
	}

	return nil
}
