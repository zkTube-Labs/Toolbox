package bsc

import (
	"fmt"
	"time"

	"github.com/nanmu42/etherscan-api"
)

const (
	// Mainnet Ethereum mainnet for production
	ScanMainnet ScanNetwork = "api"
	// Ropsten Testnet(POW)
	ScanTest ScanNetwork = "api-testnet"
)

type ScanNetwork string

func (n ScanNetwork) SubDomain() (sub string) {
	return string(n)
}

var BSCSCAN *BSCScan

type BSCScan struct {
	Cli *etherscan.Client
}

func InitBSCSCAN(env ScanNetwork, key string) {
	if BSCSCAN == nil {
		BSCSCAN = &BSCScan{
			Cli: etherscan.NewCustomized(etherscan.Customization{
				Timeout: 30 * time.Second,
				Key:     key,
				BaseURL: fmt.Sprintf(`https://%s.bscscan.com/api?`, env.SubDomain()),
			}),
		}
	}
}

func NewBSCSCAN() *BSCScan {
	return BSCSCAN
}
