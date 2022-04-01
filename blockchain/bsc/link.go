package bsc

import (
	"errors"

	"github.com/ethereum/go-ethereum/ethclient"
)

var BscRpcUrl = map[Network]string{
	Mainnet: "https://bsc-dataseed1.binance.org/",
	Test:    "https://data-seed-prebsc-1-s1.binance.org:8545/",
	Dev:     "http://127.0.0.1:8575/",
}

var BscWSUrl = map[Network]string{
	Mainnet: "wss://bsc-ws-node.nariox.org:443",
	Test:    "wss://data-seed-prebsc-1-s1.binance.org:8545/",
	Dev:     "wss://127.0.0.1:8576/",
}

var BscRpc *BSCRPC

type BSCRPC struct {
	Url   string
	WsUrl string
}

func InitBSCRPC(env Network) *BSCRPC {
	if BscRpc == nil {
		BscRpc = &BSCRPC{}
		BscRpc.SetUrl(env)
	}
	return BscRpc
}

func NewBSCRPC() *BSCRPC {
	return BscRpc
}

func (i *BSCRPC) SetUrl(k Network) (err error) {
	url, ok := BscRpcUrl[k]
	if !ok {
		err = errors.New("unsupported environment")
		return
	}
	i.Url = url
	url, ok = BscWSUrl[k]
	if !ok {
		err = errors.New("unsupported environment")
		return
	}
	i.WsUrl = url
	return
}

func (i *BSCRPC) GetClient() (*ethclient.Client, error) {
	return ethclient.Dial(i.Url)
}

func (i *BSCRPC) GetWSSClient() (*ethclient.Client, error) {
	return ethclient.Dial(i.WsUrl)
}
