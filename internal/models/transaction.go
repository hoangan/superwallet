package models

import "math/big"

type Transfer struct {
	// Unique cointID in the system
	// Different coins within the same chain
	CoinID int64   `json:"coin_id"`
	From   string  `json:"from"`
	To     string  `json:"to"`
	Value  big.Int `json:"value"`
}

type Transaction struct {
	ChainID     int64  `json:"chain_id"`
	BlockHeight int64  `json:"block_height"`
	Hash        string `json:"hash"`
	Nonce       int64  `json:"nonce"`

	// Batch transfers in single transaction
	Transfers []Transfer `json:"transfers"`
}
