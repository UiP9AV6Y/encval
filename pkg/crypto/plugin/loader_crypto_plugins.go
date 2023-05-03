//go:build crypto_plugins
// +build crypto_plugins

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

func init() {
	if LoadPath == "" {
		loadPath := []string{
			workdirLoadPath(),
			execLoadPath(),
			osLoadPath(),
		}
		LoadPath = strings.Join(loadPath, string(os.PathListSeparator))
	}

	if LoadEnv == "" {
		LoadEnv = osLoadEnv()
	}
}

// workdirLoadPath is intended for development, as plugins are built
// into the plugins directory under the worktree root
func workdirLoadPath() string {
	d, err := os.Getwd()
	if err != nil {
		d = "."
	}

	return filepath.Join(d, "plugins")
}

// execLoadPath is intended for usage with the distributable, as plugins
// are located in a directory parallel to the bin directory with the executables
func execLoadPath() string {
	d := ".."
	ex, err := os.Executable()
	if err == nil {
		d = filepath.Dir(filepath.Dir(ex))
	}

	return filepath.Join(d, "lib")
}

// osLoadPath is inteded for usage with system packages, as plugins
// are located in a system designated library location
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

func loadPlugins(dir, basename string) ([]Plugin, error) {
	ext := ".so"
	if runtime.GOOS == "windows" {
		ext = ".dll"
	}

	pat := filepath.Join(dir, basename+ext)
	matches, err := filepath.Glob(pat)
	if err != nil {
		return nil, err
	}

	result := make([]Plugin, len(matches))
	for i, m := range matches {
		p, err := loadPlugin(m)
		if err != nil {
			return nil, err
		}

		result[i] = p
	}

	return result, nil
}

func loadPlugin(file string) (Plugin, error) {
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
