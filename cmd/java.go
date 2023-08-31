package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

type javaCmd struct {
	cli *Cmd
	cmd *cobra.Command

	root            string
	contractAccount string
	fee             string
	password        string
	mnemonic        string
}

func newJavaCmd(c *Cmd) *cobra.Command {
	s := &javaCmd{
		cli: c,
	}
	s.cmd = &cobra.Command{
		Use:   "java",
		Short: "Generate fabric-sdk-java config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.handler()
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	s.addFlags()
	return s.cmd
}

func (s *javaCmd) addFlags() {
}

func (s *javaCmd) handler() error {
	return errors.New("该命令暂不支持")
}
