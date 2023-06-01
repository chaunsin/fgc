package builder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strings"

	"fgc/parse"
	"fgc/parse/host"
	"fgc/parse/mspid"

	"gopkg.in/yaml.v3"
)

func New(c host.Config, o Options) *Builder {
	msp, err := mspId.New(o.Mode)
	if err != nil {
		log.Fatalln("mspid:", err)
	}
	h, err := host.New(o.Mode, &c)
	if err != nil {
		log.Fatalln("host:", err)
	}

	t := &Builder{
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
	return t
}

func (t *Builder) Build(cc *parse.CryptoConfig) error {
	if err := cc.Valid(); err != nil {
		return fmt.Errorf("Valid:%w", err)
	}

	// client
	if err := t.client(cc); err != nil {
		return fmt.Errorf("client:%w", err)
	}

	// organizations
	if err := t.organizations(cc); err != nil {
		return fmt.Errorf("organizations:%w", err)
	}

	// channels
	if err := t.channel(cc); err != nil {
		return fmt.Errorf("channel:%w", err)
	}

	// order
	if err := t.order(cc); err != nil {
		return fmt.Errorf("order:%w", err)
	}

	// peers
	if err := t.peers(cc); err != nil {
		return fmt.Errorf("peers:%w", err)
	}

	// EntityMatchers
	if err := t.entityMatchers(cc); err != nil {
		return fmt.Errorf("entityMatchers:%w", err)
	}

	// CertificateAuthorities
	if t.opts.CA {
		if err := t.certificateAuthorities(cc); err != nil {
			return fmt.Errorf("certificateAuthorities:%w", err)
		}
	}

	// Operations
	if t.opts.CA {
		if err := t.operations(cc); err != nil {
			return fmt.Errorf("operations:%w", err)
		}
	}

	// Metrics
	if t.opts.CA {
		if err := t.metrics(cc); err != nil {
			return fmt.Errorf("metrics:%w", err)
		}
	}

	return nil
}

// Write .
func (t *Builder) Write(d []byte) (int, error) {
	return 0, nil
}

// Output .
func (t *Builder) Output(dir string) error {
	return nil
}

// JSON 生成json
func (t *Builder) JSON() ([]byte, error) {
	return json.MarshalIndent(t, "  ", " ")
}

// YAML 生成yaml
func (t *Builder) YAML() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := yaml.NewEncoder(buf)
	defer enc.Close()
	enc.SetIndent(2)
	if err := enc.Encode(t); err != nil {
		return nil, fmt.Errorf("Encode:%w", err)
	}
	return buf.Bytes(), nil
}

// client
func (t *Builder) client(cc *parse.CryptoConfig) error {
	var client = Client{
		Organization: t.opts.OrgName,
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
	if _, ok := cc.Orgs[parse.OrgName(t.opts.OrgName)]; !ok {
		for name, _ := range cc.Orgs {
			client.Organization = string(name)
			break
		}
	}

	if t.opts.Pem {
		// todo:
		// 1. 考虑可配置 golang为${FABRIC_SDK_GO_PROJECT_PATH}/${CRYPTOCONFIG_FIXTURES_PATH}
		// 2. 路径改为用户输入得路径
		client.CryptoConfig.Path = "./config/crypto-config"
	}

	// 开启双tls
	if t.opts.DoubleTls {
		var (
			clientKeyPem  PemPath
			clientCertPem PemPath
			err           error
			scp           = true
		)
		// 模糊查询对应的组织并且找到当前组织Admin的证书
		for name, org := range cc.Orgs {
			if strings.Contains(string(name), t.opts.OrgName) {
				for domain, user := range org.Users {
					if domain.UserName() == t.opts.User {
						clientKeyPem, err = newPemPath(t.opts.Pem, user.TLS.Key)
						if err != nil {
							return fmt.Errorf("newPemPath:%w", err)
						}
						clientCertPem, err = newPemPath(t.opts.Pem, user.TLS.Cert)
						if err != nil {
							return fmt.Errorf("newPemPath:%w", err)
						}
					}
				}
			}
		}

		// windows 不支持证书池
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

	t.Client = client
	return nil
}

// organizations
//#
//# list of participating organizations in this network
//#
//organizations:
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
func (t *Builder) organizations(cc *parse.CryptoConfig) error {
	for name, org := range cc.Orgs {
		var (
			peers   = make([]string, 0, len(org.Users))
			keyPem  PemPath
			certPem PemPath
			oao     OrgAndOrder
			err     error
			o       = string(name)
			//o       = name.Org()
		)

		if o == "" {
			o = string(name)
			log.Printf("[organizations] warn org name not match:%s\n", name)
		}

		if _, ok := t.Organizations[o]; ok {
			continue
		}

		for domain := range org.Server {
			peers = append(peers, string(domain))
		}
		oao.Peers = peers

		mspId, ok := t.mspId.GetMspId(string(name))
		if !ok {
			mspId = "{待替换}"
			log.Printf("[organizations] mspid not found %s\n", name)
		}
		oao.MspId = mspId

		// TODO: 待实现
		if t.opts.CA {
			oao.CertificateAuthorities = []string{}
		}

		for domain, user := range org.Users {
			if domain.UserName() == t.opts.User {
				if t.opts.Pem {
					// TODO:
					// 1. 待确认`{username}`是否是代表程序自动的去找用户,还是说需要我们自己替换用户
					// 2. cryptoPath 支持绝对路径需要考虑
					// peerOrganizations/org1.example.com/users/{username}@org1.example.com/msp
					oao.CryptoPath = fmt.Sprintf("peerOrganizations/%s/users/%s/msp", name, domain)
					continue
				}

				keyPem, err = newPemPath(t.opts.Pem, user.Msp.KeyStore.Key)
				if err != nil {
					return fmt.Errorf("newPemPath:%w", err)
				}
				certPem, err = newPemPath(t.opts.Pem, user.Msp.SignCerts.Cert)
				if err != nil {
					return fmt.Errorf("newPemPath:%w", err)
				}

				oao.Users = map[string]KC{
					t.opts.User: {
						Key:  Key{PemPath: keyPem},
						Cert: Cert{PemPath: certPem},
					},
				}
			}
		}

		t.Organizations[o] = oao
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
			log.Printf("[organizations] warn order org name not match:%s\n", name)
		}

		if _, ok := t.Organizations[o]; ok {
			continue
		}

		mspId, ok := t.mspId.GetMspId(string(name))
		if !ok {
			log.Printf("[organizations] mspid not found %s\n", name)
			mspId = "{待替换}"
		}
		oao.MspId = mspId

		for domain, user := range order.Users {
			if domain.UserName() == t.opts.User {
				if t.opts.Pem {
					oao.CryptoPath = fmt.Sprintf("peerOrganizations/%s/users/%s/msp", name, domain)
					continue
				}

				keyPem, err = newPemPath(t.opts.Pem, user.Msp.KeyStore.Key)
				if err != nil {
					return fmt.Errorf("newPemPath:%w", err)
				}
				certPem, err = newPemPath(t.opts.Pem, user.Msp.SignCerts.Cert)
				if err != nil {
					return fmt.Errorf("newPemPath:%w", err)
				}

				oao.Users = map[string]KC{
					t.opts.User: {
						Key:  Key{PemPath: keyPem},
						Cert: Cert{PemPath: certPem},
					},
				}
			}
		}
		// todo:默认使用第一条数据当有多个order？
		t.Organizations["ordererorg"] = oao
		break
	}

	return nil
}

// channel
func (t *Builder) channel(cc *parse.CryptoConfig) error {
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

	t.Channels[t.opts.ChannelName] = ChannelPeer{
		Peer:   peer,
		Policy: Policy{}, // todo: 考虑补充
		//Order:  nil, // note: golang中早期配置中可以配置order通道但是后期v1.0.0sdk配置会提示废弃因此此处不生成
	}
	return nil
}

// order
func (t *Builder) order(cc *parse.CryptoConfig) error {
	for _, order := range cc.Order {
		for domain := range order.Server {
			if _, ok := t.Orderers[string(domain)]; ok {
				continue
			}

			tlsCaCerts, err := newPemPath(t.opts.Pem, order.TLSCA.Cert)
			if err != nil {
				return fmt.Errorf("newPemPath:%s", err)
			}

			var port = "${PORT}"
			if h, ok := t.host.GetHost(string(domain)); !ok {
				log.Printf("[order] not found port:%s", domain)
			} else {
				port = h.Port()
			}

			t.Orderers[string(domain)] = Payload{
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
func (t *Builder) peers(cc *parse.CryptoConfig) error {
	for _, org := range cc.Orgs {
		for domain := range org.Server {
			if _, ok := t.Orderers[string(domain)]; ok {
				continue
			}

			tlsCaCerts, err := newPemPath(t.opts.Pem, org.TLSCA.Cert)
			if err != nil {
				return fmt.Errorf("newPemPath:%w", err)
			}

			var port = "${PORT}"
			if h, ok := t.host.GetHost(string(domain)); !ok {
				log.Printf("[order] not found port:%s", domain)
			} else {
				port = h.Port()
			}

			t.Peers[string(domain)] = Payload{
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
func (t *Builder) entityMatchers(cc *parse.CryptoConfig) error {
	var (
		peer  = make([]Matcher, 0, len(cc.Orgs)*2)
		order = make([]Matcher, 0, len(cc.Order)*2)
		ca    = make([]Matcher, 0, 2)
	)

	for _, org := range cc.Orgs {
		for domain := range org.Server {
			url, ok := t.host.GetHost(string(domain))
			if !ok {
				log.Printf("[entityMatchers] not found host:%s\n", domain)
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
			url, ok := t.host.GetHost(string(domain))
			if !ok {
				log.Printf("[entityMatchers] not found host:%s\n", domain)
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
	//for _, order := range cc.Order {
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
	//}

	t.EntityMatchers.Peer = peer
	t.EntityMatchers.Orderer = order
	t.EntityMatchers.CertificateAuthority = ca
	return nil
}

// certificateAuthorities
func (t *Builder) certificateAuthorities(cc *parse.CryptoConfig) error {

	return nil
}

// operations
func (t *Builder) operations(cc *parse.CryptoConfig) error {

	return nil
}

// metrics
func (t *Builder) metrics(cc *parse.CryptoConfig) error {

	return nil
}
