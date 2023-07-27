package host

import "context"

var (
	defaultList = map[string]Host{
		"peer0.org1.example.com": "0.0.0.0:7051",
		"peer1.org1.example.com": "0.0.0.0:8051",
		"peer0.org2.example.com": "0.0.0.0:9051",
		"peer1.org2.example.com": "0.0.0.0:10051",
		"peer0.org3.example.com": "0.0.0.0:11051",
		"peer1.org3.example.com": "0.0.0.0:12051",
		"orderer.example.com":    "0.0.0.0:7050",
		"orderer2.example.com":   "0.0.0.0:8050",
		"orderer3.example.com":   "0.0.0.0:9050",
		"orderer4.example.com":   "0.0.0.0:10050",
		"orderer5.example.com":   "0.0.0.0:11050",
	}
)

type defaultHost struct{}

func NewDefault(ctx context.Context) (FetchHost, error) {
	return &defaultHost{}, nil
}

func (d *defaultHost) GetHost(domain string) (host Host, ok bool) {
	host, ok = defaultList[domain]
	return
}

func (d *defaultHost) Close() error {
	return nil
}
