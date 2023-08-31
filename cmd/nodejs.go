package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

type nodejsCmd struct {
	cli *Cmd
	cmd *cobra.Command

	root            string
	contractAccount string
	fee             string
	password        string
	mnemonic        string
}

func newNodeJSCmd(c *Cmd) *cobra.Command {
	s := &nodejsCmd{
		cli: c,
	}
	s.cmd = &cobra.Command{
		Use:   "nodejs",
		Short: "Generate fabric-sdk-node config file",
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
}

func (s *nodejsCmd) handler() error {
	return errors.New("该命令暂不支持")
}
