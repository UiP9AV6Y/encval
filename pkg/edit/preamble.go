package edit

import (
	libio "io"
	"text/template"

	"github.com/abiosoft/lineprefix"

	"github.com/UiP9AV6Y/encval/pkg/io"
)

const (
	Prefix           = "# |"
	PreambleTemplate = `This is {{ .AppName }} edit mode. This text (lines starting with {{ .Prefix }} at the top of
the file) will be removed when you save and exit.
 - To edit encrypted values, change the content of the DEC(<num>)::{{ index .Providers 0 }}[]!
   block{{ if gt (len .Providers) 1 }} (or {{ slice .Providers 1 | join " or " }}){{ end }}.
   WARNING: DO NOT change the number in the parentheses.
 - To add a new encrypted value copy and paste a new block from the
   appropriate example below. Note that:
    * the text to encrypt goes in the square brackets
    * ensure you include the exclamation mark when you copy and paste
    * you must not include a number when adding a new block
   e.g. {{ decorate "DEC::%s[]!" .Providers | join " -or- " }}
`
)

type PreambleData struct {
	AppName   string
	Prefix    string
	Providers []string
}

func NewPreambleData() *PreambleData {
	result := &PreambleData{
		Prefix: Prefix,
	}

	return result
}

func (d *PreambleData) Write(w libio.Writer) error {
	funcMap := template.FuncMap{
		"join":     join,
		"decorate": decorate,
	}
	tmpl, err := template.New("preamble").Funcs(funcMap).Parse(PreambleTemplate)
	if err != nil {
		return err
	}
	o := lineprefix.New(
		lineprefix.Writer(w),
		lineprefix.Prefix(d.Prefix),
	)

	return tmpl.Execute(o, d)
}

// NewPreambleReader provides a io.Reader implementation
// which skips the preamble from the provided reader instance
func NewPreambleReader(r libio.Reader) *io.PrefixScanner {
	return io.NewPrefixScanner(r, []byte(Prefix))
}

// PreambleLength measures number of lines resembling a preamble.
// The result is the number of bytes required to skip in order
// to omit any preamble data.
func PreambleLength(r libio.Reader) (int64, error) {
	return io.SeekPrefixedLinesEnd(r, []byte(Prefix))
}
