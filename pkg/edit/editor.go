package edit

import (
	"errors"
	"fmt"
	libio "io"
	libos "os"
	"os/exec"

	"github.com/UiP9AV6Y/encval/pkg/io"
)

const EditorEnvVar = "EDITOR"

var (
	ErrNoEditorFound = errors.New("Unable to find editor")
	EditorPaths      = []string{
		"/usr/bin/sensible-editor",
		"/usr/bin/editor",
		"/usr/bin/vim",
		"/usr/bin/vi",
	}
)

type EditorFunc func(*libos.File) error

func FindEditorBinary() string {
	env := libos.Getenv(EditorEnvVar)
	if env != "" {
		return env
	}

	for _, p := range EditorPaths {
		info, err := libos.Stat(p)
		if err == nil && info.Mode().IsRegular() {
			return p
		}
	}

	return ""
}

type EditorOption func(*Editor)

func EditorArgs(args ...string) EditorOption {
	return func(e *Editor) {
		e.args = args
	}
}

func EditorEnv(vars ...string) EditorOption {
	return func(e *Editor) {
		e.env = vars
	}
}

func EditorStdIn(r libio.Reader) EditorOption {
	return func(e *Editor) {
		e.stdIn = r
	}
}

func EditorStdOut(w libio.Writer) EditorOption {
	return func(e *Editor) {
		e.stdOut = w
	}
}

func EditorStdin(w libio.Writer) EditorOption {
	return func(e *Editor) {
		e.stdErr = w
	}
}

func EditorPrepare(cb EditorFunc) EditorOption {
	return func(e *Editor) {
		e.prepare = cb
	}
}

func EditorSave(cb EditorFunc) EditorOption {
	return func(e *Editor) {
		e.save = cb
	}
}

func EditorUnchanged(cb EditorFunc) EditorOption {
	return func(e *Editor) {
		e.unchanged = cb
	}
}

func EditorRetry(cb func(error) bool) EditorOption {
	return func(e *Editor) {
		e.retry = cb
	}
}

type Editor struct {
	bin  string
	args []string
	env  []string

	stdIn  libio.Reader
	stdOut libio.Writer
	stdErr libio.Writer

	prepare   EditorFunc
	save      EditorFunc
	unchanged EditorFunc

	retry func(error) bool
}

func NewEditor(bin string, options ...EditorOption) (*Editor, error) {
	var err error

	if bin == "" {
		bin = FindEditorBinary()
		if bin == "" {
			return nil, ErrNoEditorFound
		}
	}

	bin, err = exec.LookPath(bin)
	if err != nil && !errors.Is(err, exec.ErrDot) {
		return nil, err
	}

	result := &Editor{
		bin:    bin,
		stdIn:  libos.Stdin,
		stdOut: libos.Stdout,
		stdErr: libos.Stderr,
	}

	for _, o := range options {
		o(result)
	}

	return result, nil
}

func (e *Editor) OpenTemp(ident string) error {
	session, err := libos.CreateTemp("", ident)
	if err != nil {
		return err
	}
	defer libos.Remove(session.Name())

	if err := e.Open(session); err != nil {
		return err
	}
	if err := session.Close(); err != nil {
		return err
	}

	return nil
}

func (e *Editor) Open(f *libos.File) (err error) {
	if e.prepare != nil {
		_, err = f.Seek(0, libio.SeekStart)
		if err != nil {
			err = fmt.Errorf("Unable to rewind file %q: %w", f.Name(), err)
			return
		}

		err = e.prepare(f)
		if err != nil {
			err = fmt.Errorf("Unable to prepare file %q: %w", f.Name(), err)
			return
		}
	}

	var retry bool
	var before, after uint32

	for {
		_, err = f.Seek(0, libio.SeekStart)
		if err != nil {
			break
		}
		before, err = io.CalculateChecksum(f)
		if err != nil {
			break
		}

		err = e.run(f)
		if err != nil {
			break
		}

		_, err = f.Seek(0, libio.SeekStart)
		if err != nil {
			break
		}
		after, err = io.CalculateChecksum(f)
		if err != nil {
			break
		}

		_, err = f.Seek(0, libio.SeekStart)
		if err != nil {
			break
		}

		if before != after {
			retry, err = e.confirm(e.save, f, "Unable to save file")
			if !retry {
				break
			}
		} else {
			retry, err = e.confirm(e.unchanged, f, "Unable to close file")
			if !retry {
				break
			}
		}
	}

	if err != nil {
		err = fmt.Errorf("Unable to edit file %q: %w", f.Name(), err)
	}

	return
}

func (e *Editor) run(f *libos.File) error {
	cmd, err := e.command(f.Name())
	if err != nil {
		return err
	}

	err = cmd.Run()
	if exErr, ok := err.(*exec.ExitError); ok {
		return fmt.Errorf("Unable to edit %q: %s", f.Name(), exErr.Stderr)
	}

	return err
}

func (e *Editor) command(file string) (*exec.Cmd, error) {
	var args []string
	if len(e.args) > 0 {
		args = make([]string, 0, len(e.args)+1)
		args = append(args, e.args...)
		args = append(args, file)
	} else {
		args = []string{file}
	}

	result := exec.Command(e.bin, args...)
	if result.Err != nil {
		return nil, result.Err
	}

	result.Env = e.env
	result.Stdin = e.stdIn
	result.Stdout = e.stdOut
	result.Stderr = e.stdErr

	return result, nil
}

func (e *Editor) confirm(proc EditorFunc, f *libos.File, message string) (bool, error) {
	if proc == nil {
		return false, nil
	}

	err := proc(f)
	if err == nil {
		return false, nil
	}

	prompt := fmt.Sprint(message+": ", err, "\nTry again?")
	rt, err2 := io.ReadConfirmation(prompt, e.stdIn, e.stdOut)

	if err2 != nil {
		return false, err2
	}

	return rt, err
}
