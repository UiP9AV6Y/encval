package cmd

import (
	"os"

	"github.com/UiP9AV6Y/encval/pkg/crypto"
	"github.com/UiP9AV6Y/encval/pkg/parser"
)

type Encrypt struct {
	Parse bool     `kong:"help='Parse the input for tokens instead of encrypting everything',negatable,short='p'"`
	File  *os.File `kong:"arg,placeholder='FILE',default='-'"`
}

func NewEncrypt() *Encrypt {
	result := &Encrypt{}

	return result
}

func (c *Encrypt) Run(ctx *GlobalOptions) error {
	defer c.File.Close()

	reg, err := ctx.NewEncrypters()
	if err != nil {
		return err
	}

	if ctx.Logger().Debug().Enabled() {
		ctx.Logger().Debug().Println("Available encryption providers: ", reg.Providers())
	}

	var tokens parser.Tokens
	var pErr error
	var data int

	if c.Parse {
		tokens, pErr = c.parseTokens()
		if pErr != nil {
			return err
		}
		ctx.Logger().Info().Printf("Found %d tokens in parsed input\n", len(tokens))
	} else {
		tokens, data, pErr = c.parseData(ctx.DefaultProvider())
		if pErr != nil {
			return err
		}
		ctx.Logger().Info().Printf("Got %d bytes as encryption input\n", data)
	}

	result, err := tokens.Convert(reg)
	if err != nil {
		return err
	}

	w, err := result.WriteTo(ctx.Writer())
	ctx.Logger().Debug().Printf("Encryption resulted in %d bytes of data\n", w)
	return err
}

func (c *Encrypt) parseTokens() (parser.Tokens, error) {
	enc, err := parser.NewDecryptionParser()
	if err != nil {
		return nil, err
	}

	return enc.ParseReader(c.File)
}

func (c *Encrypt) parseData(provider string) (parser.Tokens, int, error) {
	data, err := crypto.ParseData(c.File)
	if err != nil {
		return nil, 0, err
	}

	tokens := []parser.Token{
		parser.NewDecryptedToken(data, provider, 0, 0),
	}

	return tokens, len(data), nil
}
