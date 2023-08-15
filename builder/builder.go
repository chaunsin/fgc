package builder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/chaunsin/fgc/parse"
	"github.com/chaunsin/fgc/parse/host"
	"github.com/chaunsin/fgc/parse/mspId"

	"gopkg.in/yaml.v3"
)

func New(c host.Config, o Options) *Builder {
	msp, err := mspId.New(o.Mode)
	if err != nil {
		log.Fatalln("mspid:", err)
	}
	h, err := host.New(context.TODO(), o.Mode, &c)
	if err != nil {
		log.Fatalln("host:", err)
	}

	b := &Builder{
		opts:                   o,
		mspId:                  msp,
		host:                   h,
		Version:                "v1.0.0",
		Channels:               make(map[string]ChannelPeer, 4),
		Organizations:          make(map[string]OrgAndOrder),
		CertificateAuthorities: make(map[string]CertificateAuthorities, 1),
		Orderers:               make(map[string]Payload, 1),
		Peers:                  make(map[string]Payload, 2),
	}
	return b
}

func (b *Builder) Build(cc *parse.CryptoConfig) error {
	if err := cc.Valid(); err != nil {
		return fmt.Errorf("Valid:%w", err)
	}

	// client
	if err := b.client(cc); err != nil {
		return fmt.Errorf("client:%w", err)
	}

	// organizations
	if err := b.organizations(cc); err != nil {
		return fmt.Errorf("organizations:%w", err)
	}

	// channels
	if err := b.channel(cc); err != nil {
		return fmt.Errorf("channel:%w", err)
	}

	// order
	if err := b.order(cc); err != nil {
		return fmt.Errorf("order:%w", err)
	}

	// peers
	if err := b.peers(cc); err != nil {
		return fmt.Errorf("peers:%w", err)
	}

	// EntityMatchers
	if err := b.entityMatchers(cc); err != nil {
		return fmt.Errorf("entityMatchers:%w", err)
	}

	// CertificateAuthorities
	if b.opts.CA {
		if err := b.certificateAuthorities(cc); err != nil {
			return fmt.Errorf("certificateAuthorities:%w", err)
		}
	}

	// Operations
	if b.opts.CA {
		if err := b.operations(cc); err != nil {
			return fmt.Errorf("operations:%w", err)
		}
	}

	// Metrics
	if b.opts.CA {
		if err := b.metrics(cc); err != nil {
			return fmt.Errorf("metrics:%w", err)
		}
	}

	return nil
}

// Write .
func (b *Builder) Write(d []byte) (int, error) {
	return 0, nil
}

// Output .
func (b *Builder) Output(dir string) error {
	return nil
}

// JSON 生成json
func (b *Builder) JSON() ([]byte, error) {
	return json.MarshalIndent(b, "  ", " ")
}

// YAML 生成yaml
func (b *Builder) YAML() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := yaml.NewEncoder(buf)
	defer enc.Close()
	enc.SetIndent(2)
	if err := enc.Encode(b); err != nil {
		return nil, fmt.Errorf("Encode:%w", err)
	}
	return buf.Bytes(), nil
}

// Content 根据配置类型生成相应格式内容
func (b *Builder) Content() ([]byte, error) {
	var (
		content []byte
		err     error
	)

	switch b.opts.FileType {
	case "json":
		content, err = b.JSON()
	case "yaml":
		fallthrough
	default:
		content, err = b.YAML()
	}
	return content, err
}

// client
func (b *Builder) client(cc *parse.CryptoConfig) error {
	var client = Client{
		Organization: b.opts.OrgName,
		Logging:      Logging{Level: "info"},
		CryptoConfig: Path{Path: ""},
		CredentialStore: CredentialStore{
			Path: "./data/keystore", // todo:考虑可配置
			CryptoStore: Path{
				Path: "./data/msp", // todo:考虑可配置
			},
		},
		Bccsp: Bccsp{
			Security: Security{
				Enabled: true,
				Default: struct {
					Provider string `json:"provider,omitempty" yaml:"provider"`
				}{Provider: "SW"},
				HashAlgorithm: "SHA2",
				SoftVerify:    true,
				Level:         256,
			},
		},
	}

	// 如果输入的组织名不存在则随机取一个组织名称
	if _, ok := cc.Orgs[parse.OrgName(b.opts.OrgName)]; !ok {
		for name, _ := range cc.Orgs {
			client.Organization = string(name)
			break
		}
	}

	if b.opts.Pem {
		// todo:
		// 1. 考虑可配置 golang为${FABRIC_SDK_GO_PROJECT_PATH}/${CRYPTOCONFIG_FIXTURES_PATH}
		// 2. 路径改为用户输入得路径
		client.CryptoConfig.Path = "./config/crypto-config"
	}

	// 开启双tls
	if b.opts.DoubleTls {
		var (
			clientKeyPem  PemPath
			clientCertPem PemPath
			err           error
			scp           = true
		)
		// 模糊查询对应的组织并且找到当前组织Admin的证书
		for name, org := range cc.Orgs {
			if strings.Contains(string(name), b.opts.OrgName) {
				for domain, user := range org.Users {
					if domain.UserName() == b.opts.User {
						clientKeyPem, err = newPemPath(b.opts.Pem, user.TLS.Key)
						if err != nil {
							return fmt.Errorf("newPemPath:%w", err)
						}
						clientCertPem, err = newPemPath(b.opts.Pem, user.TLS.Cert)
						if err != nil {
							return fmt.Errorf("newPemPath:%w", err)
						}
					}
				}
			}
		}

		// windows go1.17版本不支持证书池
		if runtime.GOARCH == "Windows" {
			scp = false
		}

		client.Bccsp.TLSCerts = TLSCerts{
			SystemCertPool: scp,
			Client: KC{ // 根据组织选择证书
				Key:  Key{PemPath: clientKeyPem},
				Cert: Cert{PemPath: clientCertPem},
			},
		}
	}

	b.Client = client
	return nil
}

// organizations
// #
// # list of participating organizations in this network
// #
// organizations:
//  org1:
//    # mspid 这个值在configtx.yaml中organizations作用域下面去寻找对应
//    mspid: Org1MSP
//
//    # This org's MSP store (absolute path or relative to client.cryptoconfig)
//    cryptoPath: peerOrganizations/org1.example.com/users/{username}@org1.example.com/msp
//
//    peers:
//      - peer0.org1.example.com
//      - peer1.org1.example.com
//
//    # [Optional]. Certificate Authorities issue certificates for identification purposes in a Fabric based
//    # network. Typically certificates provisioning is done in a separate process outside of the
//    # runtime network. Fabric-CA is a special certificate authority that provides a REST APIs for
//    # dynamic certificate management (enroll, revoke, re-enroll). The following section is only for
//    # Fabric-CA servers.
//    certificateAuthorities:
//      - ca.org1.example.com
//      - tlsca.org1.example.com
//
//  # the profile will contain public information about organizations other than the one it belongs to.
//  # These are necessary information to make transaction lifecycles work, including MSP IDs and
//  # peers with a public URL to send transaction proposals. The file will not contain private
//  # information reserved for members of the organization, such as admin key and certificate,
//  # fabric-ca registrar enroll ID and secret, etc.
//  org2:
//    mspid: Org2MSP
//
//    # This org's MSP store (absolute path or relative to client.cryptoconfig)
//    cryptoPath: peerOrganizations/org2.example.com/users/{username}@org2.example.com/msp
//
//    peers:
//      - peer0.org2.example.com
//      - peer1.org2.example.com
//
//    certificateAuthorities:
//      - ca.org2.example.com
//
//  # Orderer Org name
//  ordererorg:
//    # Membership Service Provider ID for this organization
//    mspID: OrdererMSP
//
//    # Needed to load users crypto keys and certs for this org (absolute path or relative to global crypto path, DEV mode)
//    cryptoPath: ordererOrganizations/example.com/users/{username}@example.com/msp
func (b *Builder) organizations(cc *parse.CryptoConfig) error {
	for name, org := range cc.Orgs {
		var (
			peers   = make([]string, 0, len(org.Users))
			keyPem  PemPath
			certPem PemPath
			oao     OrgAndOrder
			err     error
			o       = string(name)
			// o       = name.Org()
		)

		if o == "" {
			o = string(name)
			log.Printf("[organizations] warn org name not match:%s\n", name)
		}

		if _, ok := b.Organizations[o]; ok {
			continue
		}

		for domain := range org.Server {
			peers = append(peers, string(domain))
		}
		oao.Peers = peers

		mi, ok := b.mspId.GetMspId(string(name))
		if !ok {
			mi = "{待替换}"
			log.Printf("[organizations] mspid not found: %s\n", name)
		}
		oao.MspId = mi

		// TODO: 待实现
		if b.opts.CA {
			oao.CertificateAuthorities = []string{}
		}

		for domain, user := range org.Users {
			if domain.UserName() == b.opts.User {
				if b.opts.Pem {
					// TODO:
					// 1. 待确认`{username}`是否是代表程序自动的去找用户,还是说需要我们自己替换用户
					// 2. cryptoPath 支持绝对路径需要考虑
					// peerOrganizations/org1.example.com/users/{username}@org1.example.com/msp
					oao.CryptoPath = fmt.Sprintf("peerOrganizations/%s/users/%s/msp", name, domain)
					continue
				}

				keyPem, err = newPemPath(b.opts.Pem, user.Msp.KeyStore.Key)
				if err != nil {
					return fmt.Errorf("newPemPath:%w", err)
				}
				certPem, err = newPemPath(b.opts.Pem, user.Msp.SignCerts.Cert)
				if err != nil {
					return fmt.Errorf("newPemPath:%w", err)
				}

				oao.Users = map[string]KC{
					b.opts.User: {
						Key:  Key{PemPath: keyPem},
						Cert: Cert{PemPath: certPem},
					},
				}
			}
		}

		b.Organizations[o] = oao
	}

	// TODO:ordererorg 查看官方示例中配置中带排序节点配置此处待研究,以下内容需要待验证暂时先放这里
	for name, order := range cc.Order {
		var (
			keyPem  PemPath
			certPem PemPath
			oao     OrgAndOrder
			err     error
			o       = name.Org()
		)

		if o == "" {
			o = string(name)
			log.Printf("[organizations] order org name not match: %s\n", name)
		}

		if _, ok := b.Organizations[o]; ok {
			continue
		}

		mi, ok := b.mspId.GetMspId(string(name))
		if !ok {
			mi = "{待替换}"
			log.Printf("[organizations] mspid not found: %s\n", name)
		}
		oao.MspId = mi

		for domain, user := range order.Users {
			if domain.UserName() == b.opts.User {
				if b.opts.Pem {
					oao.CryptoPath = fmt.Sprintf("peerOrganizations/%s/users/%s/msp", name, domain)
					continue
				}

				keyPem, err = newPemPath(b.opts.Pem, user.Msp.KeyStore.Key)
				if err != nil {
					return fmt.Errorf("newPemPath:%w", err)
				}
				certPem, err = newPemPath(b.opts.Pem, user.Msp.SignCerts.Cert)
				if err != nil {
					return fmt.Errorf("newPemPath:%w", err)
				}

				oao.Users = map[string]KC{
					b.opts.User: {
						Key:  Key{PemPath: keyPem},
						Cert: Cert{PemPath: certPem},
					},
				}
			}
		}
		// todo:默认使用第一条数据当有多个order？
		b.Organizations["ordererorg"] = oao
		break
	}

	return nil
}

// channel
func (b *Builder) channel(cc *parse.CryptoConfig) error {
	var peer = make(map[string]PeerPolicy, 4)
	for _, org := range cc.Orgs {
		for domain := range org.Server {
			// note: 此处用的是map序列化,由于golang的特性map无序,在实际链接中也不影响
			if _, ok := peer[string(domain)]; !ok {
				peer[string(domain)] = PeerPolicy{
					EndorsingPeer:  true,
					ChaincodeQuery: true,
					LedgerQuery:    true,
					EventSource:    true,
				}
			}
		}
	}

	b.Channels[b.opts.ChannelName] = ChannelPeer{
		Peer:   peer,
		Policy: Policy{}, // todo: 考虑补充
		// Order:  nil, // note: golang中早期配置中可以配置order通道但是后期v1.0.0sdk配置会提示废弃因此此处不生成
	}
	return nil
}

// order
func (b *Builder) order(cc *parse.CryptoConfig) error {
	for _, order := range cc.Order {
		for domain := range order.Server {
			if _, ok := b.Orderers[string(domain)]; ok {
				continue
			}

			tlsCaCerts, err := newPemPath(b.opts.Pem, order.TLSCA.Cert)
			if err != nil {
				return fmt.Errorf("newPemPath:%s", err)
			}

			var port = "${PORT}"
			if h, ok := b.host.GetHost(string(domain)); !ok {
				log.Printf("[order] not found port: %s", domain)
			} else {
				port = h.Port()
			}

			b.Orderers[string(domain)] = Payload{
				Url: fmt.Sprintf("%s:%s", domain, port),
				GrpcOptions: GrpcOptions{
					SSLTargetNameOverride: string(domain),
					KeepAliveTime:         0,
					KeepAliveTimeout:      0,
					KeepAlivePermit:       0,
					FailFast:              false,
					AllowInsecure:         false,
				},
				TlsCACerts: tlsCaCerts,
			}
		}
	}
	return nil
}

// peers
func (b *Builder) peers(cc *parse.CryptoConfig) error {
	for _, org := range cc.Orgs {
		for domain := range org.Server {
			if _, ok := b.Orderers[string(domain)]; ok {
				continue
			}

			tlsCaCerts, err := newPemPath(b.opts.Pem, org.TLSCA.Cert)
			if err != nil {
				return fmt.Errorf("newPemPath:%w", err)
			}

			var port = "${PORT}"
			if h, ok := b.host.GetHost(string(domain)); !ok {
				log.Printf("[order] not found port:%s", domain)
			} else {
				port = h.Port()
			}

			b.Peers[string(domain)] = Payload{
				Url: fmt.Sprintf("%s:%s", domain, port),
				GrpcOptions: GrpcOptions{
					SSLTargetNameOverride: string(domain),
					KeepAliveTime:         0,
					KeepAliveTimeout:      0,
					KeepAlivePermit:       0,
					FailFast:              false,
					AllowInsecure:         false,
				},
				TlsCACerts: tlsCaCerts,
			}
		}
	}
	return nil
}

// entityMatchers
func (b *Builder) entityMatchers(cc *parse.CryptoConfig) error {
	var (
		peer  = make([]Matcher, 0, len(cc.Orgs)*2)
		order = make([]Matcher, 0, len(cc.Order)*2)
		ca    = make([]Matcher, 0, 2)
	)

	for _, org := range cc.Orgs {
		for domain := range org.Server {
			url, ok := b.host.GetHost(string(domain))
			if !ok {
				log.Printf("[entityMatchers] not found host: %s\n", domain)
				url = host.Host(domain)
			}
			peer = append(peer, Matcher{
				Pattern:                             fmt.Sprintf("(\\w*)%s(\\w*)", string(domain)), // todo:考虑正则规则
				UrlSubstitutionExp:                  fmt.Sprintf("grpcs://%s:%s", url.IP(), url.Port()),
				SSLTargetOverrideUrlSubstitutionExp: string(domain),
				MappedHost:                          string(domain),
				MappedName:                          "", // todo:
				IgnoreEndpoint:                      false,
			})
		}
	}

	for _, o := range cc.Order {
		for domain := range o.Server {
			url, ok := b.host.GetHost(string(domain))
			if !ok {
				log.Printf("[entityMatchers] not found host: %s\n", domain)
				url = host.Host(domain)
			}
			order = append(order, Matcher{
				Pattern:                             fmt.Sprintf("(\\w*)%s(\\w*)", string(domain)), // todo:考虑正则规则
				UrlSubstitutionExp:                  fmt.Sprintf("grpcs://%s:%s", url.IP(), url.Port()),
				SSLTargetOverrideUrlSubstitutionExp: string(domain),
				MappedHost:                          string(domain),
				MappedName:                          "", // todo:
				IgnoreEndpoint:                      false,
			})
		}
	}

	// todo:CertificateAuthority
	// for _, order := range cc.Order {
	//	for domain := range order.Server {
	//		peer = append(peer, Matcher{
	//			Pattern:                             fmt.Sprintf("(\\w*)%s(\\w*)", string(domain)), // todo:考虑正则规则
	//			UrlSubstitutionExp:                  "grpcs://{ip}:{port}",
	//			SSLTargetOverrideUrlSubstitutionExp: string(domain),
	//			MappedHost:                          string(domain),
	//			MappedName:                          "", // todo:
	//			IgnoreEndpoint:                      false,
	//		})
	//	}
	// }

	b.EntityMatchers.Peer = peer
	b.EntityMatchers.Orderer = order
	b.EntityMatchers.CertificateAuthority = ca
	return nil
}

// certificateAuthorities
func (b *Builder) certificateAuthorities(cc *parse.CryptoConfig) error {

	return nil
}

// operations
func (b *Builder) operations(cc *parse.CryptoConfig) error {

	return nil
}

// metrics
func (b *Builder) metrics(cc *parse.CryptoConfig) error {

	return nil
}
