package main

import (
	"fmt"
	"os"
	"path/filepath"

	"filippo.io/age"

	"github.com/UiP9AV6Y/encval/pkg/fs"
)

const IdentitiesFileExt = ".key.age"

type Identities string

func (i Identities) Load() ([]age.Identity, error) {
	pattern := filepath.Join(string(i), "*"+IdentitiesFileExt)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	result := make([]age.Identity, len(matches))
	for j, m := range matches {
		data, err := os.ReadFile(m)
		if err != nil {
			return nil, err
		}

		rcpt, err := age.ParseX25519Identity(string(data))
		if err != nil {
			return nil, err
		}

		result[j] = rcpt
	}

	return result, nil
}

func (i Identities) Save(rcpt age.Identity, force bool) (string, error) {
	x25519, ok := rcpt.(*age.X25519Identity)
	if !ok {
		return "", fmt.Errorf("Unable to save unsupported identity implementation '%T'", rcpt)
	}

	path := filepath.Join(string(i), x25519.Recipient().String()+IdentitiesFileExt)

  if force {
		return path, fs.ReinstallFile(path, []byte(x25519.String()), 0644)
	}

	return path, fs.InstallFile(path, []byte(x25519.String()), 0644)
}
