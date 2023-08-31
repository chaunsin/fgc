package cmd

import (
	"fmt"
	"os"

	"github.com/chaunsin/fgc/builder"
	"github.com/chaunsin/fgc/parse"

	"github.com/spf13/cobra"
)

type golangCmd struct {
	cli *Cmd
	cmd *cobra.Command

	pathway bool // 开启golang路径魔法变量
}

func newGolangCmd(c *Cmd) *cobra.Command {
	s := &golangCmd{
		cli: c,
	}
	s.cmd = &cobra.Command{
		Use:   "go",
		Short: "golang",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.generate()
		},
	}
	s.addFlags()
	return s.cmd
}

func (s *golangCmd) addFlags() {
	s.cmd.Flags().BoolVar(&s.pathway, "pathway", false, "是否开启go的魔法变量路径")
}

func (s *golangCmd) generate() error {
	opts := s.cli.RootOpts
	opts.Language = "golang"

	// 根据模式读取文件
	cc, err := parse.Open(opts.Input, opts.Mode)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}

	// 模板对象
	b := builder.New(opts.Config, opts.Options)
	if err := b.Build(cc); err != nil {
		return fmt.Errorf("build: %w", err)
	}
	content, err := b.Content()
	if err != nil {
		return fmt.Errorf("serialize:%w", err)
	}

	if opts.Stdout {
		fmt.Fprintf(os.Stdout, "##### CONTEXT #####\n%s\n", content)
	}

	if err := os.MkdirAll(opts.Output, os.ModePerm); err != nil {
		return fmt.Errorf("output path %s invalid", opts.Output)
	}
	dir := fmt.Sprintf("%s/config.%s", opts.Output, opts.FileType)
	if err := os.WriteFile(dir, content, os.ModePerm); err != nil {
		return fmt.Errorf("WriteFile:%s", err)
	}
	return nil
}
