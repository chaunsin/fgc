package mspId

var (
	defaultMspId = map[string]string{
		"org1.example.com":     "Org1MSP",
		"org2.example.com":     "Org2MSP",
		"org3.example.com":     "Org3MSP",
		"org4.example.com":     "Org4MSP",
		"org5.example.com":     "Org5MSP",
		"orderer.example.com":  "OrdererMSP",
		"orderer2.example.com": "Orderer2MSP",
		"orderer3.example.com": "Orderer3MSP",
		"orderer4.example.com": "Orderer4MSP",
		"orderer5.example.com": "Orderer5MSP",
	}
)

type defaultMsp struct{}

func NewDefault() (FetchMspId, error) {
	return &defaultMsp{}, nil
}

func (d *defaultMsp) GetMspId(org string) (id string, ok bool) {
	id, ok = defaultMspId[org]
	return
}
