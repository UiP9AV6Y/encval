package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/UiP9AV6Y/encval/pkg/crypto"
	"github.com/UiP9AV6Y/encval/pkg/log"
)

type CreateKeys struct {
	Force   bool          `kong:"help='Overwrite existing data',negatable,short='f'"`
	Timeout time.Duration `kong:"help='Time limit for secreat generation',short='t',default='5s'"`
}

func NewCreateKeys() *CreateKeys {
	result := &CreateKeys{}

	return result
}

func (c *CreateKeys) Run(ctx *GlobalOptions) error {
	reg, err := ctx.NewEncrypters()
	if err != nil {
		return err
	}

	if ctx.Logger().Debug().Enabled() {
		ctx.Logger().Debug().Println("Available encryption providers: ", reg.Providers())
	}

	enc, ok := reg.Get(ctx.EncryptMethod)
	if !ok {
		return fmt.Errorf("%q is not a valid encryption method", ctx.EncryptMethod)
	}

	ctx.Logger().Debug().Println("Using encryption provider", ctx.EncryptMethod)

	gen, ok := enc.(crypto.SecretsGenerator)
	if !ok {
		return fmt.Errorf("%q does not require keys", ctx.EncryptMethod)
	}

	ctx.Logger().Debug().Println("Generating keys for encryption with", ctx.EncryptMethod)
	return c.generateSecrets(gen.GenerateSecrets, log.NewContext(context.Background(), ctx.Logger()))
}

func (c *CreateKeys) generateSecrets(worker func(bool, context.Context) error, ctx context.Context) error {
	ch := make(chan error, 1)
	ctxTimeout, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()

	go func(force bool) {
		ch <- worker(force, ctxTimeout)
	}(c.Force)

	select {
	case <-ctxTimeout.Done():
		return ctxTimeout.Err()
	case result := <-ch:
		return result
	}
}
