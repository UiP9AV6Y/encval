package cli

import (
	"github.com/alecthomas/kong"

	"github.com/UiP9AV6Y/encval/internal/cli/cmd"
)

type Root struct {
	*cmd.GlobalOptions

	CreateKeys *cmd.CreateKeys `kong:"cmd,help='Create a set of keys with which to encrypt/decrypt data',aliases='init,createkeys,generate-secrets'"`
	Decrypt    *cmd.Decrypt    `kong:"cmd,help='Decrypt some data',aliases='decipher,decode'"`
	Edit       *cmd.Edit       `kong:"cmd,help='Edit an encrypted file',aliases='open'"`
	Encrypt    *cmd.Encrypt    `kong:"cmd,help='Encrypt some data',aliases='encipher,encode'"`
	Password   *cmd.Password   `kong:"cmd,help='Encrypt a password entered on the terminal',aliases='pass,prompt'"`
	Recrypt    *cmd.Recrypt    `kong:"cmd,help='Recrypt some data',aliases='recipher,recode'"`
	Version    *cmd.Version    `kong:"cmd,help='Show version information'"`
}

func NewRoot(appName string, cfgs ConfigPaths, plugins kong.Plugins) *Root {
	options := cmd.NewGlobalOptions(
		cmd.GlobalAppName(appName),
		cmd.GlobalBaseDir(cfgs.Directory),
		cmd.GlobalEncryptionPlugins(plugins),
	)
	result := &Root{
		GlobalOptions: options,
		CreateKeys:    cmd.NewCreateKeys(),
		Decrypt:       cmd.NewDecrypt(),
		Edit:          cmd.NewEdit(),
		Encrypt:       cmd.NewEncrypt(),
		Password:      cmd.NewPassword(),
		Recrypt:       cmd.NewRecrypt(),
		Version:       cmd.NewVersion(),
	}

	return result
}

func (r *Root) CliGroups() []kong.Group {
	global := r.GlobalOptions.CliGroups()
	result := make([]kong.Group, 0, len(global))

	result = append(result, global...)

	return result
}
