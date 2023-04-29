package version

import (
	"runtime/debug"

	_ "github.com/UiP9AV6Y/encval/version"
)

var (
	release   = "v0.0.0"
	reference = "HEAD"
)

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			reference = s.Value
		}
	}

	needle := info.Main.Path + "/version"
	for _, d := range info.Deps {
		if d.Path == needle {
			release = d.Version
			break
		}
	}
}

func Release() string {
	return release
}

func Reference() string {
	return reference
}
