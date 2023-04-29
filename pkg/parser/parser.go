package parser

import (
	"bytes"
	"errors"
	"io"

	"github.com/timtadh/lexmachine"
	"github.com/timtadh/lexmachine/machines"

	"github.com/UiP9AV6Y/encval/pkg/crypto"
)

var (
	EncryptedPattern  = []byte(`ENC\[\w+(,[^\]]+)+\]`)
	DecryptedPattern  = []byte(`DEC(\(\d+\))?::(\w+)\[[^\]]+\]!`)
	ErrMalformedToken = errors.New("Malformed Token")
)

func ReaderBytes(r io.Reader) ([]byte, error) {
	buf := new(bytes.Buffer)

	if _, err := buf.ReadFrom(r); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type Parser struct {
	lexer   *lexmachine.Lexer
	pattern []byte
}

func NewEncryptionParser() (*Parser, error) {
	return newParser(EncryptedPattern, EncryptedTokenFactory())
}

func NewDecryptionParser() (*Parser, error) {
	return newParser(DecryptedPattern, DecryptedTokenFactory())
}

func newParser(pattern []byte, action lexmachine.Action) (*Parser, error) {
	lexer := lexmachine.NewLexer()
	lexer.Add(pattern, action)

	if err := lexer.Compile(); err != nil {
		return nil, err
	}

	result := &Parser{
		lexer:   lexer,
		pattern: pattern,
	}

	return result, nil
}

func (p *Parser) String() string {
	return string(p.pattern)
}

func (p *Parser) ParseReader(r io.Reader) (Tokens, error) {
	b, err := ReaderBytes(r)
	if err != nil {
		return nil, err
	}

	return p.ParseBytes(b)
}

func (p *Parser) ParseBytes(b []byte) (Tokens, error) {
	scanner, err := p.lexer.Scanner(b)
	if err != nil {
		return nil, err
	}

	var tokens Tokens
	buf := new(bytes.Buffer)

	for tok, err, eof := scanner.Next(); !eof; tok, err, eof = scanner.Next() {
		if err == nil {
			if buf.Len() > 0 {
				tokens = append(tokens, NewNonMatchToken(crypto.NewData(buf)))
				buf = new(bytes.Buffer)
			}

			tokens = append(tokens, tok.(Token))
		} else if data, ok := err.(*machines.UnconsumedInput); ok {
			scanner.TC = data.FailTC
			buf.Write(data.Text[data.StartTC:data.FailTC])
		} else {
			return nil, err
		}
	}

	if buf.Len() > 0 {
		tokens = append(tokens, NewNonMatchToken(crypto.NewData(buf)))
	}

	return tokens, nil
}
