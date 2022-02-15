package eth

import (
	"errors"

	"github.com/ethereum/go-ethereum/ethclient"
)

var Infura *infura

type infura struct {
	Key string
	Url string
}

var InfuraEnv = map[string]string{
	Mainnet: "mainnet",
	Rinkeby: "rinkeby",
	Ropsten: "ropsten",
}

func InitInfura(key, env string) *infura {
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

func (i *infura) SetUrl(k string) (err error) {
	env, ok := InfuraEnv[k]
	if !ok {
		err = errors.New("unsupported environment")
		return
	}
	i.Url = "https://" + env + ".infura.io/v3/" + i.Key
	return
}

func (i *infura) GetClient() (*ethclient.Client, error) {
	return ethclient.Dial(i.Url)
}
