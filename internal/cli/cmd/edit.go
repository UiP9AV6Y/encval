package cmd

import (
	"errors"
	libio "io"
	libfs "io/fs"
	libos "os"

	"github.com/UiP9AV6Y/encval/pkg/crypto"
	"github.com/UiP9AV6Y/encval/pkg/edit"
	"github.com/UiP9AV6Y/encval/pkg/parser"
)

type Edit struct {
	Preamble bool   `kong:"help='Prefix edit sessions with the informative preamble',negatable,default='true'"`
	Decrypt  bool   `kong:"help='Decrypt existing encrypted content. New content marked properly will be encrypted',negatable,default='true'"`
	Editor   string `kong:"help='Command to invoke for editing if $EDITOR is not set'"`
	File     string `kong:"arg,type='path',placeholder='FILE'"`
}

func NewEdit() *Edit {
	result := &Edit{}

	return result
}

func (c *Edit) Run(ctx *GlobalOptions) error {
	reg, err := ctx.NewEncrypters()
	if err != nil {
		return err
	}

	if ctx.Logger().Debug().Enabled() {
		ctx.Logger().Debug().Println("Available encryption providers: ", reg.Providers())
	}

	ed, err := edit.NewEditor(c.Editor,
		edit.EditorPrepare(func(f *libos.File) error {
			if c.Preamble {
				provs := reorderStringSlice(reg.Providers(), ctx.DefaultProvider())
				data := edit.NewPreambleData()
				data.AppName = ctx.AppName()
				data.Providers = provs

				ctx.Logger().Debug().Println("Writing preamble")

				if err := data.Write(f); err != nil {
					return err
				}
			}

			src, err := libos.Open(c.File)
			if err != nil {
				if errors.Is(err, libfs.ErrNotExist) {
					// there is nothing to decrypt when editing a new file
					ctx.Logger().Info().Println("Editing non-existent file", c.File)
					return nil
				}

				return err
			}
			defer src.Close()

			if c.Decrypt {
				ctx.Logger().Debug().Println("Decrypting existent file", c.File)
				return c.decrypt(f, src, reg)
			}

			return c.copy(f, src)
		}),
		edit.EditorSave(func(f *libos.File) error {
			if c.Preamble {
				ctx.Logger().Debug().Println("Removing preamble")

				if _, err := f.Seek(0, libio.SeekStart); err != nil {
					return err
				}

				skip, err := edit.PreambleLength(f)
				if err != nil {
					return err
				}

				if _, err := f.Seek(skip, libio.SeekStart); err != nil {
					return err
				}
			}

			return c.encrypt(f, reg)
		}),
		edit.EditorRetry(func(err error) bool {
			ctx.Logger().Debug().Printf("Error while editing file %q: %v\n", c.File, err)
			_, ok := err.(*parser.ParserError)
			return ok
		}),
	)
	if err != nil {
		return err
	}

	return ed.OpenTemp(ctx.AppName())
}

func (c *Edit) copy(output libio.Writer, input libio.Reader) error {
	_, err := libio.Copy(output, input)

	return err
}

func (c *Edit) decrypt(output libio.Writer, input libio.Reader, reg crypto.Encrypters) error {
	enc, err := parser.NewEncryptionParser()
	if err != nil {
		return err
	}

	tokens, err := enc.ParseReader(input)
	if err != nil {
		return err
	}

	result, err := tokens.Convert(reg)
	if err != nil {
		return err
	}

	result.Reindex()

	_, err = result.WriteTo(output)
	return err
}

func (c *Edit) encrypt(input libio.Reader, reg crypto.Encrypters) error {
	enc, err := parser.NewDecryptionParser()
	if err != nil {
		return err
	}

	tokens, err := enc.ParseReader(input)
	if err != nil {
		return err
	}

	if err := tokens.Validate(); err != nil {
		return err
	}

	result, err := tokens.Convert(reg)
	if err != nil {
		return err
	}

	output, err := libos.Create(c.File)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = result.WriteTo(output)
	return err
}

func reorderStringSlice(values []string, lead string) []string {
	result := make([]string, 0, len(values)+1)

	result = append(result, lead)

	for _, value := range values {
		if lead != value {
			result = append(result, value)
		}
	}

	return result
}
