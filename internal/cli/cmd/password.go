package cmd

import (
	"errors"
	"fmt"
	libio "io"
	libos "os"

	"golang.org/x/term"

	"github.com/UiP9AV6Y/encval/pkg/crypto"
	"github.com/UiP9AV6Y/encval/pkg/io"
	"github.com/UiP9AV6Y/encval/pkg/parser"
)

type Password struct {
}

func NewPassword() *Password {
	result := &Password{}

	return result
}

func (c *Password) Run(ctx *GlobalOptions) error {
	reg, err := ctx.NewEncrypters()
	if err != nil {
		return err
	}

	if ctx.Logger().Debug().Enabled() {
		ctx.Logger().Debug().Println("Available encryption providers: ", reg.Providers())
	}

	response, err := c.readPassword(ctx.Reader(), ctx.Writer())
	if err != nil {
		if errors.Is(err, libio.EOF) {
			// enforce newline in case of interupt.
			// we can not do this in readPassword
			// as the terminal might still be in RAW
			// mode, where a newline might be rendered
			// misaligned
			fmt.Fprintln(ctx.Writer())
			ctx.Logger().Debug().Println("Password entry was interrupted")
			return nil
		}
		return err
	}

	pwd := parser.NewDecryptedToken(response, ctx.DefaultProvider(), 0, 0)
	result, err := pwd.Convert(reg)
	if err != nil {
		return err
	}

	ctx.Logger().Debug().Println("Encrypting password using", ctx.DefaultProvider())

	n, err := fmt.Fprintln(ctx.Writer(), result)
	ctx.Logger().Debug().Printf("Encryption resulted in %d bytes of data\n", n)
	return err
}

func (c *Password) readPassword(r libio.Reader, w libio.Writer) (crypto.Data, error) {
	if f, ok := r.(*libos.File); ok && term.IsTerminal(int(f.Fd())) {
		s, err := term.MakeRaw(int(f.Fd()))
		if err != nil {
			return nil, err
		}

		defer term.Restore(int(f.Fd()), s)
	}

	rw := io.NewReadWriter(r, w)
	t := term.NewTerminal(rw, "> ")
	response, err := t.ReadPassword("Enter password: ")
	if err != nil {
		return nil, err
	}

	return crypto.Data(response), nil
}
