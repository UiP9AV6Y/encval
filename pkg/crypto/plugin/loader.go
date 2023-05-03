package plugin

import (
	"os"
	"strings"
)

var (
	// LoadPath is an OS-specific list of directories to search for plugins
	LoadPath string
	// LoadEnv is the name of the environment variable to retrieve the LoadPath
	LoadEnv string
)

// PluginLoader is a wrapper around the golang plugin system to
// create crypto.EncrypterPlugin instances
type PluginLoader string

// LoadOsPlugins searches for plugins in OS-specific locations
func (l PluginLoader) LoadOsPlugins() ([]Plugin, error) {
	fallback := strings.ReplaceAll(LoadPath, "{0}", string(l))

	return l.LoadEnvPlugins(LoadEnv, fallback)
}

// LoadEnvPlugins searches for plugins in locations designated by the given
// environment variable. The value is expected to be a list of directories.
// The item delimiter is OS-specific (os.PathListSeparator)
func (l PluginLoader) LoadEnvPlugins(key, fallback string) ([]Plugin, error) {
	path := os.Getenv(key)
	if path == "" {
		path = fallback
	}
	if path == "" {
		return []Plugin{}, nil
	}

	dirs := strings.Split(path, string(os.PathListSeparator))
	return l.LoadDirsPlugins(dirs)
}

// LoadDirsPlugins searches for plugins in the given directories
func (l PluginLoader) LoadDirsPlugins(dirs []string) ([]Plugin, error) {
	result := []Plugin{}
	for _, d := range dirs {
		p, err := l.LoadPlugins(d)
		if err != nil {
			return nil, err
		}

		result = append(result, p...)
	}

	return result, nil
}

// LoadPlugins searches for plugins in the given directory
func (l PluginLoader) LoadPlugins(dir string) ([]Plugin, error) {
	return loadPlugins(dir, string(l)+"_*")
}

// LoadPlugin uses the golang plugin system to create a Plugin instance.
func (l PluginLoader) LoadPlugin(file string) (Plugin, error) {
	return loadPlugin(file)
}
