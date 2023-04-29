package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"

	"github.com/UiP9AV6Y/encval/internal/cli/cmd"
)

const (
	AppName     = "encval"
	ConfigFile  = "config.toml"
	Description = "Inline value encryption utility"
)

func ExitIfErr(err error) {
	if err != nil {
		fmt.Println(AppName+":", err)
		os.Exit(1)
	}
}

type Cli struct {
	parser *kong.Kong
	root   *Root
}

func New() (*Cli, error) {
	plugins, err := cmd.EncryptionPlugins(AppName)
	if err != nil {
		return nil, err
	}

	cfgs := NewConfigPaths()
	root := NewRoot(AppName, cfgs, plugins)
	parser, err := kong.New(root,
		kong.Name(AppName),
		kong.DefaultEnvars(strings.ToUpper(AppName)),
		kong.Description(Description),
		kong.Configuration(TOML, cfgs...),
		kong.ShortUsageOnError(),
		kong.ExplicitGroups(root.CliGroups()),
	)
	if err != nil {
		return nil, err
	}

	result := &Cli{
		root:   root,
		parser: parser,
	}

	return result, nil
}

func (c *Cli) ExitIfErr(err error, showUsage bool) {
	if showUsage {
		c.parser.FatalIfErrorf(err)
	} else {
		ExitIfErr(err)
	}
}

func (c *Cli) Context() *cmd.GlobalOptions {
	return c.root.GlobalOptions
}

func (c *Cli) Parse() (*kong.Context, error) {
	return c.parse(os.Args[1:])
}

func (c *Cli) parse(args []string) (*kong.Context, error) {
	return c.parser.Parse(args)
}
