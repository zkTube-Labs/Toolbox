package eth

import (
	"github.com/nanmu42/etherscan-api"
)

var ETHSCAN *ETHScan

type ETHScan struct {
	cli *etherscan.Client
}

func InitETHScan(env etherscan.Network, key string) {
	if ETHSCAN == nil {
		ETHSCAN = &ETHScan{
			cli: etherscan.New(env, key),
		}
	}
}

func NewETHScan() *ETHScan {
	return ETHSCAN
}
