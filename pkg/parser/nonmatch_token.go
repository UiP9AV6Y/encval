package parser

import (
	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"

	"github.com/UiP9AV6Y/encval/pkg/crypto"
)

type NonMatchToken struct {
	value crypto.Data
}

func NonMatchTokenFactory() lexmachine.Action {
	return func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
		token := NewNonMatchToken(crypto.Data(m.Bytes))

		return token, nil
	}
}

func NewNonMatchToken(value crypto.Data) *NonMatchToken {
	result := &NonMatchToken{
		value: value,
	}

	return result
}

func (t *NonMatchToken) Provider() string {
	return ""
}

func (t *NonMatchToken) Value() crypto.Data {
	return t.value.Copy()
}

func (t *NonMatchToken) String() string {
	return string(t.value)
}

func (t *NonMatchToken) Convert(_ crypto.Encrypters) (Token, error) {
	return t, nil
}

func (t *NonMatchToken) Column() int {
	return 0
}

func (t *NonMatchToken) Line() int {
	return 0
}
