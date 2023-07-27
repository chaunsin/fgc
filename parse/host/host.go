package host

import (
	"context"
	"log"
	"strings"
)

type FetchHost interface {
	GetHost(domain string) (host Host, ok bool)
	Close() error
}

type Host string

func (d Host) Port() string {
	// 0.0.0.0:9051->9051/tcp
	if p := strings.Split(string(d), "->"); len(p) == 2 {
		if p := strings.Split(p[0], ":"); len(p) == 2 {
			return p[1]
		}
	}

	// 0.0.0.0:9051
	if p := strings.Split(string(d), ":"); len(p) == 2 {
		return p[1]
	}

	// 7051/tcp todo:

	// :::9051->9051/tcp todo:
	return "${PORT}"
}

// IP .
func (d Host) IP() string {
	// 0.0.0.0:9051->9051/tcp
	if p := strings.Split(string(d), "->"); len(p) == 2 {
		if p := strings.Split(p[0], ":"); len(p) == 2 {
			return p[0]
		}
	}

	// 0.0.0.0:9051
	if p := strings.Split(string(d), ":"); len(p) == 2 {
		return p[0]
	}

	// 7051/tcp todo:

	// :::9051->9051/tcp todo:

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

func New(ctx context.Context, mode string, cfg *Config) (FetchHost, error) {
	var (
		err  error
		resp FetchHost
	)
	switch mode {
	case "ftp":
		// todo:
	case "sftp":
		resp, err = NewSSH(ctx, cfg)
	default:
		// 执行host解析
		// resp, err = NewHostResolver(ctx)
	}
	if err != nil {
		log.Printf("[host] New:%s\n", err)
	}
	if err == nil && resp != nil {
		return resp, nil
	}

	// 从docker中尝试读取解析如果失败则使用本地默认值
	resp, err = NewLocalDocker(ctx)
	if err == nil {
		return resp, nil
	}
	return NewDefault(ctx)
}
