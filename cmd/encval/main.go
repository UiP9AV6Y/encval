package main

import (
	"github.com/UiP9AV6Y/encval/internal/cli"
)

func main() {
	state, err := cli.New()
	cli.ExitIfErr(err)

	ctx, err := state.Parse()
	state.ExitIfErr(err, true)

	err = ctx.Run(state.Context())
	state.ExitIfErr(err, false)
}
