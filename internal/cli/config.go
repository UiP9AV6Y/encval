package cli

import (
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"

	"github.com/UiP9AV6Y/encval/pkg/fs"
)

type ConfigPaths []string

func SystemConfigFile() string {
	return filepath.Join(fs.Separator, "etc", AppName, ConfigFile)
}

func UserConfigFile() string {
	base, err := os.UserConfigDir()
	if err != nil {
		base = filepath.Join("~", ".config")
	}

	return kong.ExpandPath(filepath.Join(base, AppName, ConfigFile))
}

func LocalConfigFiles() []string {
	pwd, err := os.Getwd()
	if err != nil {
		pwd = "."
	}

	base := "." + AppName
	dirs := fs.DirHierarchy(pwd)
	result := make([]string, len(dirs))
	for i, d := range dirs {
		result[i] = filepath.Join(d, base, ConfigFile)
	}

	return result
}

func NewConfigPaths() ConfigPaths {
	result := append(
		LocalConfigFiles(),
		UserConfigFile(),
		SystemConfigFile(),
	)

	return ConfigPaths(result)
}

// Directory tries to find the first existing directory in the lookup path.
// it works similar to the resolver logic behind kong.ConfigFlag, with the
// difference, that instead of an empty string, it yields the first item
// (unless the lookup path is empty, in which case an empty string is also
// returned)
func (c ConfigPaths) Directory() string {
	for _, f := range c {
		d := filepath.Dir(f)
		if _, err := os.Stat(d); err == nil {
			return d
		}
	}

	if len(c) > 0 {
		return filepath.Dir(c[0])
	}

	return ""
}
