package cmd

import "github.com/spf13/cobra"

type nodejsCmd struct {
	cli             *Cmd
	cmd             *cobra.Command

	root            string
	contractAccount string
	fee             string
	password        string
	mnemonic string
}

func newNodeJSCmd(c *Cmd) *cobra.Command {
	s := &nodejsCmd{
		cli: c,
	}
	s.cmd = &cobra.Command{
		Use:   "nodejs",
		Short: "nodejs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.handler()
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	s.addFlags()
	return s.cmd
}

func (s *nodejsCmd) addFlags() {
	s.cmd.Flags().StringVarP(&s.root, "root", "r", "./keys", "root账户配置文件位置")
	s.cmd.Flags().StringVarP(&s.mnemonic, "mnemonic", "m", "", "中文助记词,如果设置则使用当前设置的账号来创建合约账户")
	s.cmd.Flags().StringVarP(&s.contractAccount, "contractAccount", "c", "XC2222222222222222@xuper", "合约账号")
	s.cmd.Flags().StringVarP(&s.fee, "fee", "f", "999999999", "转账金额")
	s.cmd.Flags().StringVarP(&s.password, "password", "P", "123456", "公私钥文件密码")
}

func (s *nodejsCmd) handler() error {

	return nil
}
