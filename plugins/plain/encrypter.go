package main

import (
	"encoding/base64"

	"github.com/UiP9AV6Y/encval/pkg/crypto"
)

const Encrypter = "PLAIN"

// PlainEncrypter is a crypto.Encrypter implementation utilizing
// Base64 encoding as its encryption strategy. This is in no way
// secure and merely serves demontration purposes.
type PlainEncrypter struct {
	enc *base64.Encoding
}

// New creates a PlainEncrypter instance. The input parameters
// define the base64.Encoding to use for any data processing
func New(urlEncoding, padding bool) *PlainEncrypter {
	var enc *base64.Encoding

	if urlEncoding {
		if padding {
			enc = base64.URLEncoding
		} else {
			enc = base64.RawURLEncoding
		}
	} else {
		if padding {
			enc = base64.StdEncoding
		} else {
			enc = base64.RawStdEncoding
		}
	}

	result := &PlainEncrypter{
		enc: enc,
	}

	return result
}

func (p *PlainEncrypter) Encrypt(data crypto.Data) (crypto.Data, error) {
	dst := make([]byte, p.enc.EncodedLen(len(data)))
	p.enc.Encode(dst, data)

	return crypto.Data(dst), nil
}

func (p *PlainEncrypter) Decrypt(data crypto.Data) (crypto.Data, error) {
	dst := make([]byte, p.enc.DecodedLen(len(data)))
	n, err := p.enc.Decode(dst, data)
	if err != nil {
		return nil, err
	}

	return crypto.Data(dst[:n]), nil
}
