package eth

import (
	"errors"

	"github.com/ethereum/go-ethereum/ethclient"
)

var Infura *infura

type infura struct {
	Key   string
	Url   string
	WSUrl string
}

var InfuraEnv = map[Network]string{
	Mainnet: "mainnet",
	Rinkeby: "rinkeby",
	Ropsten: "ropsten",
}

func InitInfura(key string, env Network) *infura {
	if Infura == nil {
		Infura = &infura{}
		Infura.Key = key
		Infura.SetUrl(env)
	}
	return Infura
}

func NewInfura() *infura {
	return Infura
}

func (i *infura) SetUrl(k Network) (err error) {
	env, ok := InfuraEnv[k]
	if !ok {
		err = errors.New("unsupported environment")
		return
	}
	i.Url = "https://" + env + ".infura.io/v3/" + i.Key
	i.WSUrl = "wss://" + env + ".infura.io/v3/" + i.Key
	return
}

func (i *infura) GetClient() (*ethclient.Client, error) {
	return ethclient.Dial(i.Url)
}

func (i *infura) GetWSSClient() (*ethclient.Client, error) {
	return ethclient.Dial(i.WSUrl)
}
