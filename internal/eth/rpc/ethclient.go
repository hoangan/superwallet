package rpc

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/hoangan/superwallet/pkg/enccode/hexencoder"
	"github.com/hoangan/superwallet/pkg/httpclient"
)

type EthClient struct {
	client *httpclient.Client
}

func New(url string) *EthClient {
	return &EthClient{
		client: httpclient.NewHttpClient(url),
	}
}

func (c *EthClient) GetLatestBlock() (*Block, error) {
	// fetch the latest block
	responseBody, err := c.client.Post(getRequestPayload("eth_getBlockByNumber", []string{"latest", "true"}))
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block: %v", err)
	}

	block := &Block{}
	if err := json.Unmarshal(responseBody, block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block response: %v", err)
	}

	return block, nil
}

func (c *EthClient) GetBlockByNumber(blockNumber *big.Int) (*Block, error) {
	blockNumberHex := hexencoder.DecimalToHex(blockNumber)

	// fetch the block by number with detailed transactions
	responseBody, err := c.client.Post(getRequestPayload("eth_getBlockByNumber", []string{blockNumberHex, "true"}))
	if err != nil {
		return nil, fmt.Errorf("failed to get block by number: %v", err)
	}

	block := &Block{}
	if err := json.Unmarshal(responseBody, block); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block response: %v", err)
	}

	return block, nil
}

func (c *EthClient) GetTransactionByHash(txHash string) (*RawTransaction, error) {
	// fetch the transaction by hash
	responseBody, err := c.client.Post(getRequestPayload("eth_getTransactionByHash", []string{txHash}))
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by hash: %v", err)
	}

	tx := &RawTransaction{}
	if err := json.Unmarshal(responseBody, tx); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction response: %v", err)
	}

	return tx, nil
}

func (c *EthClient) TraceInternalTransaction(txHash string) ([]*RawTransaction, error) {
	// fetch the internal transactions by transaction hash
	responseBody, err := c.client.Post(getRequestPayload("trace_transaction", []string{txHash, "[\"trace\"]"}))
	if err != nil {
		return nil, fmt.Errorf("failed to get internal transactions: %v", err)
	}

	internalTxs := []*RawTransaction{}
	if err := json.Unmarshal(responseBody, &internalTxs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal internal transactions response: %v", err)
	}

	return internalTxs, nil
}

func getRequestPayload(method string, params []string) []byte {
	payload := struct {
		Jsonrpc string   `json:"jsonrpc"`
		Method  string   `json:"method"`
		Params  []string `json:"params"`
		Id      int      `json:"id"`
	}{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		Id:      1,
	}

	if payloadBytes, err := json.Marshal(payload); err != nil {
		return nil
	} else {
		return payloadBytes
	}
}
