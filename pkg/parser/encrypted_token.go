package parser

import (
	"bytes"
	"fmt"
	"unicode"

	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"

	"github.com/UiP9AV6Y/encval/pkg/crypto"
)

type EncryptedToken struct {
	column, line int
	value        crypto.Data
	provider     string
}

func EncryptedTokenFactory() lexmachine.Action {
	return func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
		return ParseEncryptedToken(m.Bytes, m.StartLine, m.StartColumn)
	}
}

func ParseEncryptedToken(match []byte, line, column int) (*EncryptedToken, error) {
	var provider string
	var value []byte

	if len(match) < 5 || match[3] != '[' { // ENC[]
		return nil, ErrMalformedToken
	}

	haystack := match[4 : len(match)-1] // ENC[ ]
	params := bytes.Split(haystack, []byte(","))

	if len(params) > 2 {
		return nil, ErrMalformedToken
	}

	if len(params) > 1 {
		value = ParseEncryptedValue(params[1])
	}

	if len(params) > 0 {
		provider = string(bytes.TrimSpace(params[0]))
	}

	result := NewEncryptedToken(crypto.Data(value), provider, line, column)
	return result, nil
}

func ParseEncryptedValue(v []byte) crypto.Data {
	result := make(crypto.Data, 0, len(v))

	for i := 0; i < len(v); i++ {
		c := v[i]
		if !unicode.IsSpace(rune(c)) {
			result = append(result, c)
		}
	}

	return result
}

func NewEncryptedToken(value crypto.Data, provider string, line, column int) *EncryptedToken {
	result := &EncryptedToken{
		line:     line,
		column:   column,
		value:    value,
		provider: provider,
	}

	return result
}

func (t *EncryptedToken) Provider() string {
	return t.provider
}

func (t *EncryptedToken) Value() crypto.Data {
	return t.value.Copy()
}

func (t *EncryptedToken) String() string {
	return fmt.Sprintf("ENC[%s,%s]", t.provider, t.value)
}

func (t *EncryptedToken) Convert(e crypto.Encrypters) (Token, error) {
	enc, ok := e.Get(t.provider)
	if !ok {
		return nil, fmt.Errorf("No such encryption provider %q available", t.provider)
	}

	value, err := enc.Decrypt(t.value)
	if err != nil {
		return nil, err
	}

	return NewDecryptedToken(value, t.provider, t.line, t.column), nil
}

func (t *EncryptedToken) Column() int {
	return t.column
}

func (t *EncryptedToken) Line() int {
	return t.line
}
