package crypto

import (
	"bytes"
	"io"
)

// Data is a functional wrapper around a byte slice
// providing various feature not found in the standard library
type Data []byte

// NewData creates a Data instance from the bytes.Buffer.Bytes()
func NewData(b *bytes.Buffer) Data {
	return Data(b.Bytes())
}

// ParseData drains the provided reader instance
// and creates a Data instance from the output
func ParseData(r io.Reader) (Data, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Data(b), nil
}

func (d Data) String() string {
	return string(d)
}

func (d Data) Copy() Data {
	dup := make(Data, len(d))

	copy(dup, d)

	return dup
}
