package main

import (
	"bytes"
	"context"
	"io"

	"filippo.io/age"

	"github.com/UiP9AV6Y/encval/pkg/crypto"
	"github.com/UiP9AV6Y/encval/pkg/log"
)

const Encrypter = "AGE"

// AgeEncrypter is a crypto.Encrypter implementation utilizing
// AGE encryption. Data is encrypted for all available recipients
// and decrypted using any available identities.
type AgeEncrypter struct {
	ids   Identities
	rcpts Recipients

	idsCache   []age.Identity
	rcptsCache []age.Recipient
}

// New creates a AgeEncrypter instance. The provided credential
// providers are queries during data processing; results are
// cached on first use.
func New(ids Identities, rcpts Recipients) *AgeEncrypter {
	result := &AgeEncrypter{
		ids:   ids,
		rcpts: rcpts,
	}

	return result
}

func (a *AgeEncrypter) Encrypt(data crypto.Data) (crypto.Data, error) {
	if a.rcptsCache == nil {
		r, err := a.rcpts.Load()
		if err != nil {
			return nil, err
		}
		a.rcptsCache = r
	}

	var buf bytes.Buffer
	w, err := age.Encrypt(&buf, a.rcptsCache...)
	if err != nil {
		return nil, err
	}
	if _, err := w.Write([]byte(data)); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	return crypto.Encode(buf.Bytes())
}

func (a *AgeEncrypter) Decrypt(data crypto.Data) (crypto.Data, error) {
	if a.idsCache == nil {
		i, err := a.ids.Load()
		if err != nil {
			return nil, err
		}
		a.idsCache = i
	}

	raw, err := crypto.Decode(data)
	if err != nil {
		return nil, err
	}

	src := bytes.NewReader([]byte(raw))
	r, err := age.Decrypt(src, a.idsCache...)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return nil, err
	}

	return crypto.Data(buf.Bytes()), nil
}

func (a *AgeEncrypter) GenerateSecrets(force bool, ctx context.Context) error {
	logger, err := log.FromContext(ctx)
	if err != nil {
		return err
	}

	key, err := age.GenerateX25519Identity()
	if err != nil {
		return err
	}

	r, err := a.rcpts.Save(key.Recipient(), force)
	if err != nil {
		return err
	}

	logger.Info().Printf("Saving %s recipient in %s\n", Encrypter, r)

	i, err := a.ids.Save(key, force)
	if err != nil {
		return err
	}

	logger.Info().Printf("Saving %s identity %s\n", Encrypter, i)

	return nil
}
