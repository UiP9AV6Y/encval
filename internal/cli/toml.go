package cli

import (
	"io"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/pelletier/go-toml/v2"
)

func TOML(r io.Reader) (kong.Resolver, error) {
	values := map[string]interface{}{}
	err := toml.NewDecoder(r).Decode(&values)
	if err != nil {
		return nil, err
	}
	var f kong.ResolverFunc = func(_ *kong.Context, _ *kong.Path, flag *kong.Flag) (interface{}, error) {
		name := strings.ReplaceAll(flag.Name, "-", "_")
		raw, ok := values[name]
		if ok {
			return raw, nil
		}
		raw = values
		for _, part := range strings.Split(name, ".") {
			if values, ok := raw.(map[string]interface{}); ok {
				raw, ok = values[part]
				if !ok {
					return nil, nil
				}
			} else {
				return nil, nil
			}
		}
		return raw, nil
	}

	return f, nil
}
