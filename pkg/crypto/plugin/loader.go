package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"reflect"
	"runtime"
	"strings"

	"github.com/UiP9AV6Y/encval/pkg/crypto"
)

var (
	// LoadPath is an OS-specific list of directories to search for plugins
	LoadPath string
	// LoadEnv is the name of the environment variable to retrieve the LoadPath
	LoadEnv string
)

func init() {
	if LoadPath == "" {
		LoadPath = osLoadPath()
	}

	if LoadEnv == "" {
		LoadEnv = osLoadEnv()
	}
}

func osLoadPath() string {
	switch runtime.GOOS {
	case "windows":
		return "C:\\ProgramFiles\\{0}\\plugins"
	default:
		return "/usr/lib/{0}"
	}
}

func osLoadEnv() string {
	switch runtime.GOOS {
	case "windows":
		return "PATH"
	case "darwin":
		return "DYLD_LIBRARY_PATH"
	default:
		return "LD_LIBRARY_PATH"
	}
}

// Plugin is a crypto.EncrypterPlugin implementation using an dynamically
// loaded implementation.
type Plugin struct {
	path string

	crypto.EncrypterPlugin `kong:"embed,group='encryption',prefix='plugin.'"`
}

// String returns the filesystem location the plugin was loaded from
func (p Plugin) String() string {
	return p.path
}

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
	ext := ".so"
	if runtime.GOOS == "windows" {
		ext = ".dll"
	}

	pat := filepath.Join(dir, string(l)+"_*"+ext)
	matches, err := filepath.Glob(pat)
	if err != nil {
		return nil, err
	}

	result := make([]Plugin, len(matches))
	for i, m := range matches {
		p, err := l.LoadPlugin(m)
		if err != nil {
			return nil, err
		}

		result[i] = p
	}

	return result, nil
}

// LoadPlugin uses the golang plugin system to create a Plugin instance.
func (l PluginLoader) LoadPlugin(file string) (Plugin, error) {
	result := Plugin{
		path: file,
	}
	p, err := plugin.Open(file)
	if err != nil {
		return result, err
	}

	sym, err := p.Lookup("NewEncrypterPlugin")
	if err != nil {
		return result, err
	}

	builder, ok := sym.(func() crypto.EncrypterPlugin)
	if !ok {
		return result, fmt.Errorf("%q is not a valid encryption plugin", file)
	}

	impl := builder()
	if isNil(impl) {
		return result, fmt.Errorf("%q failed to initialize", file)
	}

	result.EncrypterPlugin = impl
	return result, nil
}

// https://mangatmodi.medium.com/go-check-nil-interface-the-right-way-d142776edef1#26f4
func isNil(i interface{}) bool {
	if i == nil {
		return true
	}

	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}

	return false
}
