package cmd

import (
	"fmt"

	"github.com/UiP9AV6Y/encval/internal/version"
)

type Version struct {
}

func NewVersion() *Version {
	result := &Version{}

	return result
}

func (c *Version) Run(ctx *GlobalOptions) error {
	fmt.Println(version.Release(), version.Reference())
	return nil
}
