package fs

import (
	"os"
	"path/filepath"
)

const (
	UserReadable  os.FileMode = 0400
	GroupReadable os.FileMode = 0040
	OtherReadable os.FileMode = 0004
)

func InstallFile(name string, data []byte, perm os.FileMode) error {
	return writeFile(name, data, os.O_WRONLY|os.O_CREATE|os.O_EXCL, perm)
}

func ReinstallFile(name string, data []byte, perm os.FileMode) error {
	return writeFile(name, data, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
}

func dirMode(fileMode os.FileMode) (dirMode os.FileMode) {
	if fileMode&UserReadable != 0 {
		dirMode |= 0700
	}
	if fileMode&GroupReadable != 0 {
		dirMode |= 0050
	}
	if fileMode&OtherReadable != 0 {
		dirMode |= 0005
	}

	return
}

func writeFile(name string, data []byte, flag int, perm os.FileMode) error {
	dir := filepath.Dir(name)
	dperm := dirMode(perm)
	if err := os.MkdirAll(dir, dperm); err != nil {
		return err
	}

	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}

	return err
}
