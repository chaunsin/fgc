package mspId

type FetchMspId interface {
	GetMspId(org string) (id string, ok bool)
}

func New(mode string) (FetchMspId, error) {
	// var (
	//	err  error
	//	resp FetchMspId
	// )
	// switch mode {
	// case "docker":
	// case "ssh":
	// }
	// if err == nil && resp!=nil{
	//	return resp, nil
	// }
	return NewDefault()
}
