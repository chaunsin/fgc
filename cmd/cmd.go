package cmd

import (
	"os"

	"github.com/chaunsin/fabric-gen-config/builder"
	"github.com/chaunsin/fabric-gen-config/parse/host"

	"github.com/spf13/cobra"
)

type RootOpts struct {
	Debug    bool   // 是否开启命令行debug模式
	Input    string // 加载证书路径
	Output   string // 生成文件路径
	Stdout   bool   // 生成内容是否打印到标准数据中
	Service  string // 生成链接服务的配置类型 normal:传统方式(默认) gateway:网关方式
	FileType string // 生成文件类型 yaml(默认) json
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
			Example: "fgc golang",
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
	c.root.PersistentFlags().StringVarP(&c.RootOpts.Output, "output", "p", "./", "")
	c.root.PersistentFlags().BoolVar(&c.RootOpts.Stdout, "stdout", false, "")
	c.root.PersistentFlags().StringVarP(&c.RootOpts.FileType, "type", "t", "yaml", "")
	c.root.PersistentFlags().StringVarP(&c.RootOpts.Service, "service", "s", "normal", "")
	c.root.PersistentFlags().BoolVar(&c.RootOpts.Pem, "pem", false, "")
	c.root.PersistentFlags().BoolVar(&c.RootOpts.DoubleTls, "tls", false, "")
	c.root.PersistentFlags().BoolVar(&c.RootOpts.CA, "ca", false, "")
	c.root.PersistentFlags().BoolVar(&c.RootOpts.Metrics, "metrics", false, "")
	c.root.PersistentFlags().BoolVar(&c.RootOpts.Operations, "operations", false, "")
	c.root.PersistentFlags().StringVarP(&c.RootOpts.OrgName, "org", "o", "org1", "")
	c.root.PersistentFlags().StringVarP(&c.RootOpts.OrderName, "order", "O", "order", "")
	c.root.PersistentFlags().StringVarP(&c.RootOpts.ChannelName, "channel", "c", "mychannel", "")
	c.root.PersistentFlags().StringVarP(&c.RootOpts.User, "user", "u", "Admin", "")
	c.root.PersistentFlags().StringVarP(&c.RootOpts.Mode, "mode", "m", "local", "local,sftp,ftp")
	c.root.PersistentFlags().StringVarP(&c.RootOpts.Addr, "host", "H", "", "")
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
