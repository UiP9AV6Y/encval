//go:build !crypto_plugins
// +build !crypto_plugins

package plugin

import (
	"errors"
)

var (
	errNoop = errors.New("Plugin feature has been disabled at build time")
)

func loadPlugins(dir, basename string) ([]Plugin, error) {
	return []Plugin{}, nil
}

func loadPlugin(file string) (Plugin, error) {
	result := Plugin{
		path: file,
	}

	return result, errNoop
}
