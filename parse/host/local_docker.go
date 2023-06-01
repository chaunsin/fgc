package host

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"time"
)

const cmdArgs = `docker ps --format "table{{.Image}}\t{{.Names}}\t{{.Ports}}" | grep "hyperledger/fabric-peer\|hyperledger/fabric-orderer" | awk '{print $2,$3}'`

type docker struct {
	store map[string]string
}

// NewLocal 考虑使用docker api来实现此功能
func NewLocal() (FetchHost, error) {
	var (
		dk     docker
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", "-c", cmdArgs)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.Println("[NewLocal] GetMspId:", stderr.String())
		return nil, fmt.Errorf("Run:%s", stderr.String())
	}
	resp := stdout.String()
	log.Printf("[NewLocal] stdout:\n%s\n", resp)
	dk.store = StrToMap(resp, " ")
	if len(dk.store) == 0 {
		return nil, errors.New("is empty")
	}
	return &dk, nil
}

func (d *docker) GetHost(domain string) (host Host, ok bool) {
	m, ok := d.store[domain]
	host = Host(m)
	return
}

func (d *docker) Close() error {
	return nil
}
