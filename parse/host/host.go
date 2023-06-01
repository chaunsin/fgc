package host

import (
	"log"
	"strings"
)

type FetchHost interface {
	GetHost(domain string) (host Host, ok bool)
	Close() error
}

type Host string

func (d Host) Port() string {
	if p := strings.Split(string(d), ":"); len(p) == 2 {
		return p[1]
	}
	return "${PORT}"
}

// IP
// TODO: IP需要二次解析
// 1. 7051/tcp, 0.0.0.0:9051->9051/tcp, :::9051->9051/tcp
// 2. 0.0.0.0:7050->7050/tcp, :::7050->7050/tcp
func (d Host) IP() string {
	if p := strings.Split(string(d), ":"); len(p) == 2 {
		return p[0]
	}
	return "${IP}"
}

func StrToMap(raw string, sep string) map[string]string {
	var resp = make(map[string]string)
	re := strings.Replace(raw, `"`, "", -1)
	for _, v := range strings.Split(re, "\n") {
		if v == "" || v == " " {
			continue
		}
		s := strings.SplitN(v, sep, 2)
		resp[s[0]] = strings.TrimSpace(s[1])
	}
	return resp
}

func New(mode string, cfg *Config) (FetchHost, error) {
	var (
		err  error
		resp FetchHost
	)
	switch mode {
	case "ftp":
		// 执行docker 或者k8s等
	case "sftp":
		// 执行docker 或者k8s等
		resp, err = NewSSH(cfg)
	}
	if err != nil {
		log.Printf("[host] New:%s\n", err)
	}
	if err == nil && resp != nil {
		return resp, nil
	}
	// 从docker中读取如果失败则讲解为本地默认值
	resp, err = NewLocal()
	if err == nil {
		return resp, nil
	}
	return NewDefault()
}
