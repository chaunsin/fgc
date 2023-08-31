package cmd

import (
	"os"

	"github.com/chaunsin/fgc/builder"
	"github.com/chaunsin/fgc/parse/host"

	"github.com/spf13/cobra"
)

type RootOpts struct {
	Debug   bool   // 是否开启命令行debug模式
	Input   string // 加载证书路径
	Output  string // 生成文件路径
	Stdout  bool   // 生成内容是否打印到标准数据中
	Service string // 生成链接服务的配置类型 normal:传统方式(默认) gateway:网关方式
	builder.Options
	host.Config
}

type Cmd struct {
	root     *cobra.Command
	RootOpts RootOpts
}

func New() *Cmd {
	c := &Cmd{
		root: &cobra.Command{
			Use:     "fgc",
			Example: "fgc golang -i ./crypto-config",
		},
	}
	c.addFlags()
	c.Add(newGolangCmd(c))
	c.Add(newNodeJSCmd(c))
	c.Add(newJavaCmd(c))
	return c
}

func (c *Cmd) addFlags() {
	c.root.PersistentFlags().BoolVar(&c.RootOpts.Debug, "debug", false, "")
	c.root.PersistentFlags().StringVarP(&c.RootOpts.Input, "input", "i", defaultString("FABRIC_CFG_PATH", "./crypto-config"), "gen [command] -i ./crypto-config")
	c.root.PersistentFlags().StringVarP(&c.RootOpts.Output, "output", "p", "./", "Generate file directory location")
	c.root.PersistentFlags().BoolVar(&c.RootOpts.Stdout, "stdout", false, "")
	c.root.PersistentFlags().StringVarP(&c.RootOpts.FileType, "type", "t", "yaml", "Generated file type")
	c.root.PersistentFlags().StringVarP(&c.RootOpts.Service, "service", "s", "normal", "")
	c.root.PersistentFlags().BoolVar(&c.RootOpts.Pem, "pem", false, "")
	c.root.PersistentFlags().BoolVar(&c.RootOpts.DoubleTls, "tls", false, "Whether to enable bidirectional TLS authentication. The default value is unidirectional")
	c.root.PersistentFlags().BoolVar(&c.RootOpts.CA, "ca", false, "")
	c.root.PersistentFlags().BoolVar(&c.RootOpts.Metrics, "metrics", false, "")
	c.root.PersistentFlags().BoolVar(&c.RootOpts.Operations, "operations", false, "")
	c.root.PersistentFlags().StringVarP(&c.RootOpts.OrgName, "org", "o", "org1", "Organization name")
	c.root.PersistentFlags().StringVarP(&c.RootOpts.OrderName, "order", "O", "order", "Orderer name")
	c.root.PersistentFlags().StringVarP(&c.RootOpts.ChannelName, "channel", "c", "mychannel", "The name of the channel used")
	c.root.PersistentFlags().StringVarP(&c.RootOpts.User, "user", "u", "Admin", "The user name used")
	c.root.PersistentFlags().StringVarP(&c.RootOpts.Mode, "mode", "m", "local", "local,sftp,ftp")
	c.root.PersistentFlags().StringVarP(&c.RootOpts.Addr, "host", "H", "", "Service ip address or domain name")
	c.root.PersistentFlags().StringVarP(&c.RootOpts.Username, "username", "U", "root", "")
	c.root.PersistentFlags().StringVarP(&c.RootOpts.Password, "password", "P", "", "")
	c.root.PersistentFlags().StringVarP(&c.RootOpts.PrivateKey, "key", "k", ".ssh/key.pem", "")
}

func (c *Cmd) Version(version string) {
	c.root.Version = version
}

func (c *Cmd) Add(command ...*cobra.Command) {
	c.root.AddCommand(command...)
}

func (c *Cmd) Execute() {
	if err := c.root.Execute(); err != nil {
		panic(err)
	}
}

func defaultString(env, value string) string {
	v := os.Getenv(env)
	if v == "" {
		return value
	}
	return v
}
