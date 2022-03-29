package bsc

import eth "github.com/zkTube-Labs/Toolbox/blockchain/ethereum"

// Network is ethereum network type (mainnet, ropsten, etc)
type Network string

const (
	Mainnet Network = "mainnet"
	Test    Network = "test"
)

var Bsc *eth.Ethereum

func InitBSC(Hot, Cold string) (err error) {
	Bsc = &eth.Ethereum{
		HotWallet: Hot,
		ColdWallt: Cold,
	}
	b := NewBSCRPC()
	Bsc.Cli, err = b.GetClient()
	// if err != nil {
	// 	return
	// }
	// Bsc.WSCli, err = b.GetWSSClient()
	return
}

func NewBSC() *eth.Ethereum {
	return Bsc
}
