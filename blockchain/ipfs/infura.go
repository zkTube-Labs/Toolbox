package ipfs

import (
	"net/http"

	ipfsApi "github.com/ipfs/go-ipfs-api"
)

const InfuraHost = "https://ipfs.infura.io:5001"

var Inf *Infura

type Infura struct {
	Shell *ipfsApi.Shell
}

func InitInfura(ProjectId, ProjectSecret string) {
	if Inf != nil {
		return
	}
	Inf = &Infura{
		Shell: ipfsApi.NewShellWithClient(InfuraHost, &http.Client{
			Transport: authTransport{
				RoundTripper:  http.DefaultTransport,
				ProjectId:     ProjectId,
				ProjectSecret: ProjectSecret,
			},
		}),
	}
}

func NewInfura() *Infura {
	return Inf
}

// authTransport decorates each request with a basic auth header.
type authTransport struct {
	http.RoundTripper
	ProjectId     string
	ProjectSecret string
}

func (t authTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.SetBasicAuth(t.ProjectId, t.ProjectSecret)
	return t.RoundTripper.RoundTrip(r)
}
