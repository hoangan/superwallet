package hexencoder

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

func HexToDecimal(hexStr string) (*big.Int, error) {
	if !strings.HasPrefix(hexStr, "0x") {
		return nil, errors.New("hex string should have 0x prefix")
	}

	hexStr = hexStr[2:]

	decimalValue := new(big.Int)
	decimalValue, ok := decimalValue.SetString(hexStr, 16)
	if !ok {
		return nil, fmt.Errorf("failed to convert hex to decimal: %s", hexStr)
	}
	return decimalValue, nil
}

func DecimalToHex(decimal *big.Int) string {
	return fmt.Sprintf("0x%x", decimal)
}
