package parser

import (
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/UiP9AV6Y/encval/pkg/crypto"
)

var ErrDuplicateIndex = errors.New("A duplicate decrypted token was found based on its index identifier")

type Token interface {
	Provider() string
	Value() crypto.Data
	Convert(crypto.Encrypters) (Token, error)
	Column() int
	Line() int
}

type ParserError struct {
	Column, Line int
	Err          error
}

func NewTokenParserError(err error, t Token) *ParserError {
	result := &ParserError{
		Column: t.Column(),
		Line:   t.Line(),
		Err:    err,
	}

	return result
}

func NewFileParserError(err error, line, column int) *ParserError {
	result := &ParserError{
		Column: column,
		Line:   line,
		Err:    err,
	}

	return result
}

func (e *ParserError) Error() string {
	return fmt.Sprintf("%s at line %d column %d", e.Err.Error(), e.Line, e.Column)
}

func (e *ParserError) Unwrap() error { return e.Err }

type Tokens []Token

func (t Tokens) Convert(enc crypto.Encrypters) (Tokens, error) {
	result := make(Tokens, 0, len(t))
	for _, e := range t {
		v, err := e.Convert(enc)
		if err != nil {
			return nil, NewTokenParserError(err, e)
		}

		result = append(result, v)
	}

	return result, nil
}

func (t Tokens) WriteTo(output io.Writer) (int64, error) {
	var result int64
	for _, e := range t {
		n, err := io.WriteString(output, e.(fmt.Stringer).String())
		result = result + int64(n)
		if err != nil {
			return result, NewTokenParserError(err, e)
		}
	}

	return result, nil
}

func (t Tokens) Validate() error {
	decs := make([]*DecryptedToken, 0, len(t))

	for _, e := range t {
		d, ok := e.(*DecryptedToken)
		if ok && d.Index > 0 {
			decs = append(decs, d)
		}
	}

	sort.Slice(decs, func(i, j int) bool {
		return decs[i].Index < decs[j].Index
	})

	for i := 1; i < len(decs); i++ {
		if decs[i].Index == decs[i-1].Index {
			return NewTokenParserError(ErrDuplicateIndex, decs[i])
		}
	}

	return nil
}

func (t Tokens) Reindex() (count int) {
	for i, e := range t {
		d, ok := e.(*DecryptedToken)
		if ok {
			d.Index = i
			count++
		}
	}

	return
}

func (t Tokens) Provision(provider string) (count int) {
	for _, e := range t {
		d, ok := e.(*DecryptedToken)
		if ok {
			d.provider = provider
			count++
		}
	}

	return
}

func (t Tokens) Plaintext() (count int) {
	for i, e := range t {
		d, ok := e.(*DecryptedToken)
		if ok {
			t[i] = NewNonMatchToken(d.value)
			count++
		}
	}

	return
}
