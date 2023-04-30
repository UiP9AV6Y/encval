package cmd

import (
	"os"

	"github.com/UiP9AV6Y/encval/pkg/parser"
)

type Decrypt struct {
	Parse bool     `kong:"help='Parse the input for tokens instead of decrypting everything',negatable,short='p'"`
	File  *os.File `kong:"arg,placeholder='FILE',default='-'"`
}

func NewDecrypt() *Decrypt {
	result := &Decrypt{}

	return result
}

func (c *Decrypt) Run(ctx *GlobalOptions) error {
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

	ctx.Logger().Info().Printf("Found %d tokens in parsed input\n", len(tokens))

	result, err := tokens.Convert(reg)
	if err != nil {
		return err
	}

	if !c.Parse {
		affected := result.Plaintext()
		ctx.Logger().Debug().Printf("Decryption affected %d tokens\n", affected)
	}

	n, err := result.WriteTo(ctx.Writer())
	ctx.Logger().Debug().Printf("Decryption resulted in %d bytes of data\n", n)
	return err
}
