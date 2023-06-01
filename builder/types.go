package builder

import (
	"fmt"
	"time"

	"fgc/parse"
	"fgc/parse/host"
	"fgc/parse/mspid"
)

type Path struct {
	Path string `json:"path,omitempty" yaml:"path"`
}

type Logging struct {
	Level string `json:"level,omitempty" yaml:"level"`
}

type CredentialStore struct {
	Path        string `json:"path,omitempty" yaml:"path,omitempty"`
	CryptoStore Path   `json:"cryptoStore,omitempty" yaml:"cryptoStore,omitempty"`
}

type Security struct {
	Enabled bool `json:"enabled,omitempty" yaml:"enabled"`
	Default struct {
		Provider string `json:"provider,omitempty" yaml:"provider"`
	} `json:"default,omitempty" yaml:"default"`
	HashAlgorithm string `json:"hash_algorithm,omitempty" yaml:"hashAlgorithm"`
	SoftVerify    bool   `json:"soft_verify,omitempty" yaml:"softVerify"`
	Level         int64  `json:"level,omitempty" yaml:"level"`
}

type Key struct {
	PemPath `yaml:",inline"`
}

type Cert struct {
	PemPath `yaml:",inline"`
}

type KC struct {
	Key  Key  `json:"key,omitempty" yaml:"key"`
	Cert Cert `json:"cert,omitempty" yaml:"cert"`
}

type TLSCerts struct {
	SystemCertPool bool `json:"system_cert_pool,omitempty" yaml:"systemCertPool"`
	Client         KC   `json:"client,omitempty" yaml:"client"`
}

type Bccsp struct {
	Security Security `json:"security" yaml:"security"`
	TLSCerts TLSCerts `json:"tls_certs,omitempty" yaml:"tlsCerts,omitempty"`
}

type PemPath struct {
	Path         string `json:"path,omitempty" yaml:"path,omitempty"`
	Pem          string `json:"pem,omitempty" yaml:"pem,omitempty"`
	projectPath  string //${FABRIC_SDK_GO_PROJECT_PATH}
	fixturesPath string //${CRYPTOCONFIG_FIXTURES_PATH}
	pem          bool
}

// newPemPath isPem is true mean use path value
func newPemPath(isPem bool, f parse.File) (PemPath, error) {
	p := PemPath{
		projectPath:  "",
		fixturesPath: "",
		pem:          isPem,
	}
	if isPem {
		p.Path = f.Path()
	} else {
		pem, err := f.Open()
		if err != nil {
			return PemPath{}, fmt.Errorf("newPemPath Open:%s", err)
		}
		p.Pem = pem
	}
	return p, nil
}

func (p PemPath) ToGolangPath() string {
	if p.Path != "" {
		return fmt.Sprintf("%s/%s/%s", p.projectPath, p.fixturesPath, p.Path)
	}
	return p.Path
}

type PeerPolicy struct {
	EndorsingPeer  bool `json:"endorsing_peer" yaml:"endorsingPeer"`
	ChaincodeQuery bool `json:"chaincode_query" yaml:"chaincodeQuery"`
	LedgerQuery    bool `json:"ledger_query" yaml:"ledgerQuery"`
	EventSource    bool `json:"event_source" yaml:"eventSource"`
}

type RetryOpts struct {
	Attempts       int           `json:"attempts,omitempty" yaml:"attempts"`
	InitialBackoff time.Duration `json:"initialBackoff,omitempty" yaml:"initialBackoff"`
	MaxBackoff     time.Duration `json:"maxBackoff,omitempty" yaml:"maxBackoff"`
	BackoffFactor  float64       `json:"backoffFactor,omitempty" yaml:"backoffFactor"`
}

type Discovery struct {
	MaxTargets int       `json:"maxTargets,omitempty" yaml:"maxTargets"`
	RetryOpts  RetryOpts `json:"retryOpts,omitempty" yaml:"retryOpts"`
}

type Selection struct {
	SortingStrategy         string `json:"sortingStrategy,omitempty" yaml:"SortingStrategy"`
	Balancer                string `json:"balancer,omitempty" yaml:"Balancer"`
	BlockHeightLagThreshold int    `json:"blockHeightLagThreshold,omitempty" yaml:"BlockHeightLagThreshold"`
}

type QueryChannelConfig struct {
	MinResponses int       `json:"minResponses,omitempty" yaml:"minResponses"`
	MaxTargets   int       `json:"maxTargets,omitempty" yaml:"maxTargets"`
	RetryOpts    RetryOpts `json:"retryOpts" yaml:"retryOpts"`
}

type EventService struct {
	ResolverStrategy                 string        `json:"resolverStrategy,omitempty" yaml:"resolverStrategy"`
	Balancer                         string        `json:"balancer,omitempty" yaml:"balancer"`
	BlockHeightLagThreshold          int           `json:"blockHeightLagThreshold,omitempty" yaml:"blockHeightLagThreshold"`
	ReconnectBlockHeightLagThreshold int           `json:"reconnectBlockHeightLagThreshold,omitempty" yaml:"reconnectBlockHeightLagThreshold"`
	PeerMonitorPeriod                time.Duration `json:"peerMonitorPeriod,omitempty" yaml:"peerMonitorPeriod"`
}

type Policy struct {
	Discovery          Discovery          `json:"discovery,omitempty" yaml:"discovery"`
	Selection          Selection          `json:"selection,omitempty" yaml:"selection"`
	QueryChannelConfig QueryChannelConfig `json:"queryChannelConfig,omitempty" yaml:"queryChannelConfig"`
	EventService       EventService       `json:"eventService,omitempty" yaml:"eventService"`
}

type ChannelPeer struct {
	Peer   map[string]PeerPolicy `json:"peer,omitempty" yaml:"peers,omitempty"`    // key为组织名称
	Order  []string              `json:"order,omitempty" yaml:"order,omitempty"`   // todo: golang中早期配置中可以配置order通道但是后期v1.0.0sdk配置会提示废弃因此此处不生成
	Policy Policy                `json:"policy,omitempty" yaml:"policy,omitempty"` // 策略
}

type OrgAndOrder struct {
	MspId                  string        `json:"msp_id,omitempty" yaml:"mspid"`
	CryptoPath             string        `json:"crypto_path,omitempty" yaml:"cryptoPath,omitempty"`
	Peers                  []string      `json:"peers,omitempty" yaml:"peers,omitempty"`
	CertificateAuthorities []string      `json:"certificate_authorities,omitempty" yaml:"certificateAuthorities,omitempty"`
	Users                  map[string]KC `json:"users" yaml:"users,omitempty"`
}

type GrpcOptions struct {
	SSLTargetNameOverride string        `json:"ssl_target_name_override,omitempty" yaml:"ssl-target-name-override"`
	AllowInsecure         bool          `json:"allow_insecure,omitempty" yaml:"allow-insecure"`
	FailFast              bool          `json:"fail_fast,omitempty" yaml:"fail-fast"`
	KeepAliveTime         time.Duration `json:"keep_alive_time,omitempty" yaml:"keep-alive-time"`
	KeepAliveTimeout      time.Duration `json:"keep_alive_timeout,omitempty" yaml:"keep-alive-timeout"`
	KeepAlivePermit       time.Duration `json:"keep_alive_permit" yaml:"keep-alive-permit"`
}

type Payload struct {
	Url         string      `json:"url,omitempty" yaml:"url"`
	GrpcOptions GrpcOptions `json:"grpcOptions,omitempty" yaml:"grpcOptions"`
	TlsCACerts  PemPath     `json:"tls_ca_certs,omitempty" yaml:"tlsCACerts"`
}

type CertificateAuthoritiesTLSCACerts struct {
	PemPath `yaml:"path,inline"`
	Client  KC `yaml:"client"`
}

type Registrar struct {
	EnrollId     string `json:"enroll_id,omitempty" yaml:"enrollId"`
	EnrollSecret string `json:"enroll_secret,omitempty" yaml:"enroll_secret"`
}

type CertificateAuthorities struct {
	Url         string      `json:"url,omitempty" yaml:"url"`
	GrpcOptions GrpcOptions `json:"grpcOptions,omitempty" yaml:"grpcOptions"`
	TlsCACerts  PemPath     `json:"tls_ca_certs,omitempty" yaml:"tlsCACerts"`
	Registrar   Registrar   `json:"registrar,omitempty" yaml:"registrar"`
	CaName      string      `json:"ca_name,omitempty" yaml:"caName"`
}

type Matcher struct {
	Pattern                             string `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	UrlSubstitutionExp                  string `json:"url_substitution_exp,omitempty" yaml:"urlSubstitutionExp,omitempty"`
	SSLTargetOverrideUrlSubstitutionExp string `json:"ssl_target_override_url_substitution_exp,omitempty" yaml:"sslTargetOverrideUrlSubstitutionExp,omitempty"`
	MappedHost                          string `json:"mapped_host,omitempty" yaml:"mappedHost,omitempty"`
	MappedName                          string `json:"mapped_name,omitempty" yaml:"mappedName,omitempty"`         // TODO:读golang源码发现
	IgnoreEndpoint                      bool   `json:"ignore_endpoint,omitempty" yaml:"ignoreEndpoint,omitempty"` // TODO:读golang源码发现
}

type EntityMatchers struct {
	Peer                 []Matcher `json:"peer,omitempty" yaml:"peer,omitempty"`
	Orderer              []Matcher `json:"orderer,omitempty" yaml:"orderer,omitempty"`
	CertificateAuthority []Matcher `json:"certificate_authority" yaml:"certificateAuthority,omitempty"`
}

type OperationsTLS struct {
	enabled bool
	Cert    struct {
		File string `json:"file,omitempty" yaml:"file"`
	} `json:"cert,omitempty" yaml:"cert"`
	Key struct {
		File string `json:"file,omitempty" yaml:"file"`
	} `json:"key,omitempty" yaml:"key"`
	ClientAuthRequired bool `json:"client_auth_required,omitempty" yaml:"clientAuthRequired"`
	ClientRootCAs      struct {
		Files []string `json:"files,omitempty" yaml:"files"`
	}
}

type Operations struct {
	ListenAddress string        `json:"listen_address,omitempty" yaml:"listenAddress"`
	Tls           OperationsTLS `json:"tls,omitempty" yaml:"tls"`
}

type Statsd struct {
	Network       string        `json:"network,omitempty" yaml:"network"`
	Address       string        `json:"address,omitempty" yaml:"address"`
	WriteInterval time.Duration `json:"write_interval,omitempty" yaml:"writeInterval"`
	Prefix        string        `json:"prefix,omitempty" yaml:"prefix"`
}

type Metrics struct {
	Provider string `json:"provider,omitempty" yaml:"provider"`
	Statsd   Statsd `json:"statsd,omitempty" yaml:"statsd"`
}

type Client struct {
	Organization    string          `json:"organization,omitempty" yaml:"organization"`
	Logging         Logging         `json:"logging,omitempty" yaml:"logging"`
	CryptoConfig    Path            `json:"cryptoconfig,omitempty" yaml:"cryptoconfig,omitempty"`
	CredentialStore CredentialStore `json:"credentialStore,omitempty" yaml:"credentialStore"`
	Bccsp           Bccsp           `json:"BCCSP,omitempty" yaml:"BCCSP"`
}

type Builder struct {
	opts  Options
	host  host.FetchHost
	mspId mspId.FetchMspId

	Version                string                            `json:"version,omitempty" yaml:"version"`
	Client                 Client                            `json:"client,omitempty" yaml:"client,omitempty"`
	Organizations          map[string]OrgAndOrder            `json:"organizations,omitempty" yaml:"organizations,omitempty"` // key为组织域名
	Channels               map[string]ChannelPeer            `json:"channels,omitempty" yaml:"channels,omitempty"`           // key为通道名称
	Orderers               map[string]Payload                `json:"orderers,omitempty" yaml:"orderers,omitempty"`           // key为组织域名
	Peers                  map[string]Payload                `json:"peers,omitempty" yaml:"peers,omitempty"`                 // key为组织域名
	CertificateAuthorities map[string]CertificateAuthorities `json:"certificate_authorities,omitempty" yaml:"certificateAuthorities,omitempty"`
	EntityMatchers         EntityMatchers                    `json:"entity_matchers,omitempty" yaml:"entityMatchers,omitempty"`
	Operations             Operations                        `json:"operations,omitempty" yaml:"operations,omitempty"`
	Metrics                Metrics                           `json:"metrics,omitempty" yaml:"metrics,omitempty"`
}
