package cmd

import (
	"os"
	"strconv"
	"strings"

	"github.com/alecthomas/kong"

	"github.com/UiP9AV6Y/encval/pkg/crypto/plugin"
)

func EncryptionPlugins(appName string) (kong.Plugins, error) {
	var plugins []plugin.Plugin
	disableLoad := disableEncryptionPluginsLoading(strings.ToUpper(appName + "_disable_plugins_loading"))
	disableList := encryptionPluginsFilter(strings.ToUpper(appName + "_disabled_plugins"))

	if disableLoad {
		return kong.Plugins{
			plugin.NewDefaultPlugin(),
		}, nil
	}

	loader := plugin.PluginLoader(appName)
	plugins, err := loader.LoadOsPlugins()
	if err != nil {
		return nil, err
	}

	result := make(kong.Plugins, 0, len(plugins)+1)
	result = append(result, plugin.NewDefaultPlugin())
	for _, p := range plugins {
		if !stringSliceContains(disableList, p.Encrypter()) {
			result = append(result, p)
		}
	}

	return result, nil
}

func disableEncryptionPluginsLoading(envKey string) (ok bool) {
	v := os.Getenv(envKey)
	if v == "" {
		return
	}

	ok, _ = strconv.ParseBool(v)

	return
}

func encryptionPluginsFilter(envKey string) []string {
	v := os.Getenv(envKey)

	return strings.Split(v, ",")
}

func stringSliceContains(haystack []string, needle string) bool {
	if len(haystack) == 0 {
		return false
	}
	if len(haystack) == 1 {
		return haystack[0] == needle
	}

	for _, item := range haystack {
		if item == needle {
			return true
		}
	}

	return false
}
