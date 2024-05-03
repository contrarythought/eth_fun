package helper

import (
	"errors"
	"math"
	"math/big"
)

func EthValue(str string) (*big.Float, error) {
	fBalance := new(big.Float)
	if _, ok := fBalance.SetString(str); !ok {
		return nil, errors.New("err: input error in EthValue()")
	}

	ethValue := new(big.Float).Quo(fBalance, big.NewFloat(math.Pow10(18)))

	return ethValue, nil
}
