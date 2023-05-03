package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"

	"github.com/UiP9AV6Y/encval/pkg/crypto"
	"github.com/UiP9AV6Y/encval/pkg/log"
)

type GlobalOptions struct {
	kong.Plugins

	EncryptMethod string          `kong:"help='Override default encryption and decryption method',default='pkcs7',short='n',name='encrypt-method'"`
	Verbose       int             `kong:"help='Increase output verbosity. Repeat for increased verbosity',type='counter',short='v'"`
	Quiet         bool            `kong:"help='Disable any log output.',negatable,short='q'"`
	Config        kong.ConfigFlag `kong:"help='Path to the configuration file',placeholder='FILE'"`

	appName string               `kong:"-"`
	baseDir func() string        `kong:"-"`
	reader  io.Reader            `kong:"-"`
	writer  io.Writer            `kong:"-"`
	logger  log.LoggerController `kong:"-"`
}

type GlobalOption func(*GlobalOptions)

func GlobalBaseDir(provider func() string) GlobalOption {
	return func(o *GlobalOptions) {
		o.baseDir = provider
	}
}

func GlobalAppName(appName string) GlobalOption {
	return func(o *GlobalOptions) {
		o.appName = appName
	}
}

func GlobalReader(reader io.Reader) GlobalOption {
	return func(o *GlobalOptions) {
		o.reader = reader
	}
}

func GlobalWriter(writer io.Writer) GlobalOption {
	return func(o *GlobalOptions) {
		o.writer = writer
	}
}

func GlobalLogger(logger io.Writer) GlobalOption {
	return func(o *GlobalOptions) {
		o.logger.SetOutput(logger)
	}
}

func GlobalEncryptionPlugins(plugins kong.Plugins) GlobalOption {
	return func(o *GlobalOptions) {
		o.Plugins = append(o.Plugins, plugins...)
	}
}

func NewGlobalOptions(options ...GlobalOption) *GlobalOptions {
	baseDir := func() string {
		dir, _ := os.Getwd()
		return dir
	}
	result := &GlobalOptions{
		baseDir: baseDir,
		reader:  os.Stdin,
		writer:  os.Stdout,
		logger:  log.NewStreamLogger(os.Stderr),
	}

	for _, o := range options {
		o(result)
	}

	return result
}

func (o *GlobalOptions) AppName() string {
	return o.appName
}

func (o *GlobalOptions) Reader() io.Reader {
	return o.reader
}

func (o *GlobalOptions) Writer() io.Writer {
	return o.writer
}

func (o *GlobalOptions) Logger() log.Logger {
	return o.logger
}

func (o *GlobalOptions) DefaultProvider() string {
	return strings.ToUpper(o.EncryptMethod)
}

func (o *GlobalOptions) CliGroups() []kong.Group {
	result := []kong.Group{
		kong.Group{
			Key:         "encryption",
			Title:       "Encryption",
			Description: `Cryptography related control settings`,
		},
	}

	return result
}

func (o *GlobalOptions) AfterApply() error {
	var v log.Verbosity
	if !o.Quiet {
		// we want INFO to be the lowest level
		// configurable with the verbosity controls
		v = log.Verbosity(o.Verbose + 3)
	}

	o.logger.SetVerbosity(v)

	o.logger.Debug().Println("EncryptMethod:", o.EncryptMethod)
	o.logger.Debug().Println("Verbose:", o.Verbose)
	o.logger.Debug().Println("Quiet:", o.Quiet)
	o.logger.Debug().Println("Config:", o.Config)
	o.logger.Debug().Println("Plugins:", len(o.Plugins))

	return nil
}

func (o *GlobalOptions) NewEncrypters() (crypto.Encrypters, error) {
	var dir string
	if o.Config != "" {
		dir = filepath.Dir(string(o.Config))
	} else {
		dir = o.baseDir()
	}

	o.logger.Debug().Println("Encryption config base directory:", dir)

	result := make(crypto.Encrypters, len(o.Plugins))
	for _, p := range o.Plugins {
		b, ok := p.(crypto.EncrypterPlugin)
		if !ok {
			return nil, fmt.Errorf("Registered plugin '%T' is not an encryption plugin builder", p)
		}

		e, err := b.NewEncrypter(dir)
		if err != nil {
			return nil, err
		}

		ok = result.Add(b.Encrypter(), e)
		if !ok {
			return nil, fmt.Errorf("Unable to register %q encrypter", b.Encrypter())
		}
	}

	return result, nil
}
