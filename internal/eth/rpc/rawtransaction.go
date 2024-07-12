package rpc

type Block struct {
	Difficulty      string `json:"difficulty"`
	Extradata       string `json:"extraData"`
	Gaslimit        string `json:"gasLimit"`
	Gasused         string `json:"gasUsed"`
	Hash            string `json:"hash"`
	Logsbloom       string `json:"logsBloom"`
	Miner           string `json:"miner"`
	Mixhash         string `json:"mixHash"`
	Nonce           string `json:"nonce"`
	Number          string `json:"number"`
	Parenthash      string `json:"parentHash"`
	Receiptsroot    string `json:"receiptsRoot"`
	Sha3uncles      string `json:"sha3Uncles"`
	Size            string `json:"size"`
	Stateroot       string `json:"stateRoot"`
	Timestamp       string `json:"timestamp"`
	Totaldifficulty string `json:"totalDifficulty"`
	// Transactions
	Transactions     []*RawTransaction `json:"transactions"`
	Transactionsroot string            `json:"transactionsRoot"`
}

type RawTransaction struct {
	Type                 string `json:"type"`
	Blockhash            string `json:"blockHash"`
	Blocknumber          string `json:"blockNumber"`
	From                 string `json:"from"`
	Gas                  string `json:"gas"`
	Hash                 string `json:"hash"`
	Input                string `json:"input"`
	Nonce                string `json:"nonce"`
	To                   string `json:"to"`
	Transactionindex     string `json:"transactionIndex"`
	Value                string `json:"value"`
	V                    string `json:"v"`
	R                    string `json:"r"`
	S                    string `json:"s"`
	Gasprice             string `json:"gasPrice"`
	Maxfeepergas         string `json:"maxFeePerGas"`
	Maxpriorityfeepergas string `json:"maxPriorityFeePerGas"`
	Chainid              string `json:"chainId"`
}

// InternalTransactionDetail is for transaction with internal transfer
type InternalTransactionDetail struct{}
