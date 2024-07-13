package testdata

import (
	"github.com/hoangan/superwallet/internal/eth/rpc"
)

var (
	RawTransaction1 = &rpc.RawTransaction{
		Hash:        "0x6523b98c957773ece31e36dfe4309df7cd6cd697f70b7f4df6d39fc008cc693a",
		Type:        "0x2",
		BlockHash:   "0xf20326ecb02332687c918de6df6c8b354ccdf8406ea1b276a4da07e22b072715",
		BlockNumber: "0x1359a3b",
		ChainId:     "0x1",
		Nonce:       "0x0",
		Gas:         "0x5208",
		GasPrice:    "0x3b9aca00",
		From:        "0x4838b106fce9647bdf1e7877bf73ce8b0bad5f97",
		To:          "0x29182006a4967e9a50c0a66076da514993d3b4d4",
		Value:       "0xa588ee0d2314c0",
	}
)
