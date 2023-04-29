package main

import (
	"fmt"
	"os"
	"path/filepath"

	"filippo.io/age"

	"github.com/UiP9AV6Y/encval/pkg/fs"
)

const RecipientsFileExt = ".cert.age"

type Recipients string

func (r Recipients) Load() ([]age.Recipient, error) {
	pattern := filepath.Join(string(r), "*"+RecipientsFileExt)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	result := make([]age.Recipient, len(matches))
	for i, m := range matches {
		data, err := os.ReadFile(m)
		if err != nil {
			return nil, err
		}

		rcpt, err := age.ParseX25519Recipient(string(data))
		if err != nil {
			return nil, err
		}

		result[i] = rcpt
	}

	return result, nil
}

func (r Recipients) Save(rcpt age.Recipient, force bool) (string, error) {
	x25519, ok := rcpt.(*age.X25519Recipient)
	if !ok {
		return "", fmt.Errorf("Unable to save unsupported recipient implementation '%T'", rcpt)
	}

	path := filepath.Join(string(r), x25519.String()+RecipientsFileExt)

  if force {
		return path, fs.ReinstallFile(path, []byte(x25519.String()), 0644)
	}

	return path, fs.InstallFile(path, []byte(x25519.String()), 0644)
}
