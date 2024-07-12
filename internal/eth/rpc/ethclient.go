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

func NewEthClient(url string) *EthClient {
	return &EthClient{
		client: httpclient.NewHttpClient(url),
	}
}

func (c *EthClient) GetLatestBlock() (*RawBlock, error) {
	// fetch the latest block
	responseBodyBytes, err := c.client.Post(getRequestPayload("eth_getBlockByNumber", []interface{}{"latest", true}))
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block: %v", err)
	}

	var responseBody struct {
		Block   RawBlock `json:"result"`
		Jsonrpc string   `json:"jsonrpc"`
		Id      int      `json:"id"`
	}
	if err := json.Unmarshal(responseBodyBytes, &responseBody); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block response: %v", err)
	}

	return &responseBody.Block, nil
}

func (c *EthClient) GetBlockByNumber(blockNumber *big.Int) (*RawBlock, error) {
	blockNumberHex := hexencoder.DecimalToHex(blockNumber)

	// fetch the block by number with detailed transactions
	responseBodyBytes, err := c.client.Post(getRequestPayload("eth_getBlockByNumber", []interface{}{blockNumberHex, true}))
	if err != nil {
		return nil, fmt.Errorf("failed to get block by number: %v", err)
	}

	var responseBody struct {
		Block   RawBlock `json:"result"`
		Jsonrpc string   `json:"jsonrpc"`
		Id      int      `json:"id"`
	}
	if err := json.Unmarshal(responseBodyBytes, &responseBody); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block response: %v", err)
	}

	return &responseBody.Block, nil
}

func (c *EthClient) GetTransactionByHash(txHash string) (*RawTransaction, error) {
	// fetch the transaction by hash
	responseBodyBytes, err := c.client.Post(getRequestPayload("eth_getTransactionByHash", []interface{}{txHash}))
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction by hash: %v", err)
	}

	var responseBody struct {
		Transaction RawTransaction `json:"result"`
		Jsonrpc     string         `json:"jsonrpc"`
		Id          int            `json:"id"`
	}

	if err := json.Unmarshal(responseBodyBytes, &responseBody); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block response: %v", err)
	}

	return &responseBody.Transaction, nil
}

// func (c *EthClient) TraceInternalTransaction(txHash string) ([]*RawTransaction, error) {
// 	// fetch the internal transactions by transaction hash
// 	responseBody, err := c.client.Post(getRequestPayload("trace_transaction", []inter{txHash, "[\"trace\"]"}))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get internal transactions: %v", err)
// 	}

// 	internalTxs := []*RawTransaction{}
// 	if err := json.Unmarshal(responseBody, &internalTxs); err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal internal transactions response: %v", err)
// 	}

// 	return internalTxs, nil
// }

func getRequestPayload(method string, params []interface{}) []byte {
	payload := struct {
		Jsonrpc string        `json:"jsonrpc"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
		Id      int           `json:"id"`
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
