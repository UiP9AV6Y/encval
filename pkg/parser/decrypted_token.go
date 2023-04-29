package parser

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"

	"github.com/UiP9AV6Y/encval/pkg/crypto"
)

type DecryptedToken struct {
	column, line, Index int
	value               crypto.Data
	provider            string
}

func DecryptedTokenFactory() lexmachine.Action {
	return func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
		return ParseDecryptedToken(m.Bytes, m.StartLine, m.StartColumn)
	}
}

func ParseDecryptedToken(match []byte, line, column int) (*DecryptedToken, error) {
	// DEC(id)::ALG[]!
	// DEC::ALG[]!
	var provider string
	var value []byte
	var id int

	bracket := bytes.IndexByte(match, '[')
	parentheses := 2 // DEC

	if bracket < 0 || len(match) < 8 || match[2] != 'C' { // DEC::[]!
		return nil, ErrMalformedToken
	}

	if match[3] == '(' {
		parentheses = bytes.IndexByte(match, ')')

		if parentheses < 0 || parentheses < 4 { // DEC)( || DEC()
			return nil, ErrMalformedToken
		}

		if i, err := strconv.ParseInt(string(match[4:parentheses]), 10, 0); err != nil {
			return nil, err
		} else {
			id = int(i)
		}
	}

	if id < 0 {
		return nil, ErrMalformedToken
	}

	if parentheses > (len(match)-4) || parentheses > bracket { // DEC() || DEC([])
		return nil, ErrMalformedToken
	}

	if match[parentheses+2] != ':' {
		return nil, ErrMalformedToken
	}

	provider = string(match[parentheses+3 : bracket])
	value = match[bracket+1 : len(match)-2]

	result := NewDecryptedToken(crypto.Data(value), provider, line, column)
	result.Index = id
	return result, nil
}

func NewDecryptedToken(value crypto.Data, provider string, line, column int) *DecryptedToken {
	result := &DecryptedToken{
		line:     line,
		column:   column,
		value:    value,
		provider: provider,
	}

	return result
}

func (t *DecryptedToken) Provider() string {
	return t.provider
}

func (t *DecryptedToken) Value() crypto.Data {
	return t.value.Copy()
}

func (t *DecryptedToken) String() string {
	if t.Index > 0 {
		return fmt.Sprintf("DEC(%d)::%s[%s]!", t.Index, t.provider, t.value)
	}

	return fmt.Sprintf("DEC::%s[%s]!", t.provider, t.value)
}

func (t *DecryptedToken) Convert(e crypto.Encrypters) (Token, error) {
	enc, ok := e.Get(t.provider)
	if !ok {
		return nil, fmt.Errorf("No such encryption provider %q available", t.provider)
	}

	value, err := enc.Encrypt(t.value)
	if err != nil {
		return nil, err
	}

	return NewEncryptedToken(value, t.provider, t.line, t.column), nil
}

func (t *DecryptedToken) Column() int {
	return t.column
}

func (t *DecryptedToken) Line() int {
	return t.line
}
