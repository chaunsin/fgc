package host

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"golang.org/x/crypto/ssh"
)

type Config struct {
	Addr       string `json:"addr,omitempty" yaml:"addr"`               // ip加端口
	Username   string `json:"user,omitempty" yaml:"user"`               // 用户名
	Password   string `json:"password,omitempty" yaml:"password"`       // 密码
	PrivateKey string `json:"private_key,omitempty" yaml:"private_key"` // eg: /home/user/.ssh/id_rsa"
	Gssapi     string `json:"gssapi,omitempty" yaml:"gssapi"`           //
}

func (c *Config) Valid() error {
	if c.Username == "" {
		c.Username = "root"
	}
	if c.Addr == "" {
		return fmt.Errorf("addr is empty")
	}
	if c.Password == "" && c.PrivateKey == "" {
		return fmt.Errorf("auth is empty")
	}
	return nil
}

type SSH struct {
	*ssh.Client
	store map[string]string
}

func NewSSH(ctx context.Context, cfg *Config) (*SSH, error) {
	log.Printf("ssh config:%+v\n", cfg)
	if err := cfg.Valid(); err != nil {
		return nil, fmt.Errorf("Valid:%w", err)
	}
	var (
		// pubKey ssh.PublicKey
		auth []ssh.AuthMethod
	)
	if cfg.Password != "" {
		auth = append(auth, ssh.RetryableAuthMethod(ssh.Password(cfg.Password), 1))
	}
	if cfg.PrivateKey != "" {
		key, err := ioutil.ReadFile(cfg.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("ReadFile:%w", err)
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("ParsePrivateKey:%w", err)
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}
	if cfg.Gssapi != "" {
		// auth = append(auth, ssh.GSSAPIWithMICAuthMethod())
	}

	conf := &ssh.ClientConfig{
		User:            cfg.Username,                // 连接登录的用户
		Auth:            auth,                        // 认证方式
		BannerCallback:  ssh.BannerDisplayStderr(),   // 错误显示到标准错误终端,可以自定义方法实现比如把错误信息写到文件中
		Timeout:         time.Second * 15,            // 建立超时时间,0为永不超时
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 貌似是限制指定秘钥类型
		// HostKeyCallback:   ssh.FixedHostKey(pubKey),  // 貌似是限制指定秘钥类型
	}
	conn, err := ssh.Dial("tcp", cfg.Addr, conf)
	if err != nil {
		// return nil, fmt.Errorf("Dial:%w", err)
		log.Fatalf("Dial:%s", err)
	}

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	// 创建一个session用于执行命令行,相当于exec中的cmd命令行
	session, err := conn.NewSession()
	if err != nil {
		return nil, fmt.Errorf("NewSession:%w", err)
	}
	defer session.Close()
	session.Stdout = &stdout
	session.Stderr = &stderr
	if err := session.Run(cmdArgs); err != nil {
		log.Printf("Run stderr:%s\n", stderr.String())
		return nil, fmt.Errorf("Run:%s", stderr.String())
	}
	resp := stdout.String()
	log.Printf("[NewSSH] stdout:\n%s\n", resp)
	store := StrToMap(resp, " ")
	if len(store) == 0 {
		return nil, errors.New("is empty")
	}
	s := &SSH{
		Client: conn,
		store:  store,
	}
	return s, nil
}

func (s *SSH) GetHost(domain string) (host Host, ok bool) {
	h, ok := s.store[domain]
	host = Host(h)
	return
}

func (s *SSH) Close() error {
	return s.Client.Close()
}
