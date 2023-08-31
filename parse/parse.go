package parse

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	orderPath = "ordererOrganizations"
	peerPath  = "peerOrganizations"
)

type (
	OrgDomain string
	OrgName   string
)

// Org 尝试解析返回组织名称 TODO: 貌似可以废弃
func (o OrgName) Org() string {
	if string(o) == "" {
		return ""
	}
	for _, v := range strings.Split(string(o), ".") {
		if !strings.Contains(v, "org") {
			continue
		}
		// todo:使用正则匹配 org开头数字结尾的表达式
		if v == "org1" || v == "org2" || v == "org3" || v == "org4" || v == "org5" {
			return v
		}
		// peer0-org1-suyuan-com.org1-suyuan-com.svc.cluster.local
		if strings.HasPrefix(v, "org") && strings.HasSuffix(v, "com") {
			return v
		}
	}
	return ""
}

type UserDomain string

func (u UserDomain) UserName() string {
	if s := strings.Split(string(u), "@"); len(s) > 0 {
		return s[0]
	}
	return ""
}

type File string

func (f File) Path() string { return string(f) }

func (f File) Open() (string, error) {
	data, err := os.ReadFile(string(f))
	if err != nil {
		return "", fmt.Errorf("ReadFile:%w", err)
	}
	return string(data), nil
}

type Package struct {
	parentDir  string
	currentDir string

	CA   File // ca
	Cert File // cert.pem
	Key  File // sk
	Yaml File // config.yaml
}

type Msp struct {
	AdminCerts Package
	CaCerts    Package
	KeyStore   Package
	SignCerts  Package
	TLSCaCerts Package
	ConfigYaml Package
}

type Serve struct {
	Msp Msp
	TLS Package
}

type User struct {
	Name string // Admin User1 ...
	Msp  Msp
	TLS  Package
}

type Org struct {
	CA     Package
	Msp    Msp
	Server map[OrgDomain]*Serve // peer/order
	TLSCA  Package
	Users  map[UserDomain]*User
}

type CryptoConfig struct {
	Orgs  map[OrgName]*Org
	Order map[OrgName]*Org
}

// Valid 检验
func (c *CryptoConfig) Valid() error {
	if c == nil {
		return errors.New("CryptoConfig is nil")
	}
	if len(c.Orgs) <= 0 {
		return errors.New("orgs is empty")
	}
	if len(c.Order) <= 0 {
		return errors.New("order is empty")
	}
	return nil
}

// GetOrgName 获取peer组织名称列表
func (c *CryptoConfig) GetOrgName() []string {
	var list = make([]string, 0, len(c.Orgs))
	for name := range c.Orgs {
		list = append(list, string(name))
	}
	return list
}

// GetOrderName 获取order组织名称列表
func (c *CryptoConfig) GetOrderName() []string {
	var list = make([]string, 0, len(c.Order))
	for name := range c.Order {
		list = append(list, string(name))
	}
	return list
}

// GetOrgUser 根据peer组织名称查询有哪些用户
func (c *CryptoConfig) GetOrgUser(org string) []string {
	var list = make([]string, 0, len(c.Orgs))
	for domain := range c.Orgs[OrgName(org)].Users {
		list = append(list, domain.UserName())
	}
	return list
}

// GetOrderUser 根据order组织名称查询有哪些用户
func (c *CryptoConfig) GetOrderUser(org string) []string {
	var list = make([]string, 0, len(c.Orgs))
	for domain := range c.Order[OrgName(org)].Users {
		list = append(list, domain.UserName())
	}
	return list
}

func Open(dir string, mode string) (*CryptoConfig, error) {
	var cc = CryptoConfig{
		Orgs:  make(map[OrgName]*Org),
		Order: make(map[OrgName]*Org),
	}
	switch mode {
	case "ftp":
		fallthrough
	case "sftp":
		fallthrough
	default:
		if dir == "./crypto-config" {
			wd, err := os.Getwd()
			if err != nil {
				return nil, fmt.Errorf("getwd: %w", err)
			}
			dir = filepath.Join(wd, dir)
		}

		if err := readPeer(dir, &cc); err != nil {
			return nil, fmt.Errorf("readPeer:%w", err)
		}
		if err := readOrder(dir, &cc); err != nil {
			return nil, fmt.Errorf("readOrder:%w", err)
		}
	}

	return &cc, cc.Valid()
}
