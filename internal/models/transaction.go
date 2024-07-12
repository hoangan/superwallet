package models

import "math/big"

type Transfer struct {
	// Unique cointID in the system
	// Different coins within the same chain
	CoinID int64    `json:"coinId"`
	Ticker string   `json:"ticker"`
	From   string   `json:"from"`
	To     string   `json:"to"`
	Value  *big.Int `json:"value"`
}

type Transaction struct {
	Type             *big.Int `json:"type"`
	BlockHash        string   `json:"blockHash"`
	BlockNumber      *big.Int `json:"blockNumber"`
	From             string   `json:"from"`
	Gas              *big.Int `json:"gas"`
	Hash             string   `json:"hash"`
	Input            string   `json:"input"`
	Nonce            *big.Int `json:"nonce"`
	To               string   `json:"to"`
	ChainId          *big.Int `json:"chainId"`
	TransactionIndex *big.Int `json:"transactionIndex"`
	Value            *big.Int `json:"value"`
	GasPrice         *big.Int `json:"gasPrice"`

	// Batch transfers of coins in single transaction
	// Any values transferred recorded here
	Transfers []*Transfer `json:"transfers"`
}
