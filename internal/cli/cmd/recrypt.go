package cmd

import (
	"os"

	"github.com/UiP9AV6Y/encval/pkg/parser"
)

type Recrypt struct {
	File *os.File `kong:"arg,placeholder='FILE',default='-'"`
}

func NewRecrypt() *Recrypt {
	result := &Recrypt{}

	return result
}

func (c *Recrypt) Run(ctx *GlobalOptions) error {
	defer c.File.Close()

	reg, err := ctx.NewEncrypters()
	if err != nil {
		return err
	}

	if ctx.Logger().Debug().Enabled() {
		ctx.Logger().Debug().Println("Available encryption providers: ", reg.Providers())
	}

	enc, err := parser.NewEncryptionParser()
	if err != nil {
		return err
	}

	ctx.Logger().Debug().Println("Using encryption parser", enc)

	tokens, err := enc.ParseReader(c.File)
	if err != nil {
		return err
	}

	ctx.Logger().Debug().Printf("Found %d tokens in parsed input\n", len(tokens))

	decrypted, err := tokens.Convert(reg)
	if err != nil {
		return err
	}

	affected := decrypted.Provision(ctx.EncryptMethod)
	result, err := decrypted.Convert(reg)
	if err != nil {
		return err
	}

	ctx.Logger().Info().Printf("Re-encrypted %d tokens in parsed input\n", affected)

	n, err := result.WriteTo(ctx.Writer())
	ctx.Logger().Debug().Printf("Decryption resulted in %d bytes of data\n", n)
	return err
}
