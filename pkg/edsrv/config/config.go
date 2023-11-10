package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

// Config represents the configuration informed via command-line flags.
type Config struct {
	LogLevel *slog.Level // log verbosity level
	Addr     string      // listen address
	TmpDir   string      // temporary directory
	Editor   string      // command-line editor
}

const (
	// LogLevelFlag log-level flag name.
	LogLevelFlag = "log-level"
	// AddrFlag listen address ("addr") flag name.
	AddrFlag = "addr"
	// TmpDirFlag temporary directory ("tmp-dir") flag name.
	TmpDirFlag = "tmp-dir"
	// EditorFlag editor command and args ("editor") flag name.
	EditorFlag = "editor"
)

// ErrInvalidConfig shows the configuration is invalid, missing elements.
var ErrInvalidConfig = errors.New("invalid configuration")

// AddLogLevelFlag adds the "log-level" flag to configure its verbosity level.
func (c *Config) AddLogLevelFlag(f *pflag.FlagSet) {
	f.Var(
		NewLogLevelValue(c.LogLevel),
		LogLevelFlag,
		fmt.Sprintf(
			"log verbosity level (default %q)",
			strings.ToLower(c.LogLevel.String()),
		),
	)
}

// AddAddrFlag adds "addr" flag, also present on "status" subcommand.
func (c *Config) AddAddrFlag(f *pflag.FlagSet) {
	f.StringVar(&c.Addr, AddrFlag, c.Addr, "listen address and port")
}

// AddTmpDirFlag adds "tmp-dir" flag.
func (c *Config) AddTmpDirFlag(f *pflag.FlagSet) {
	f.StringVar(&c.TmpDir, TmpDirFlag, c.TmpDir, "temporary directory")
}

// AddEditorFlag adds "editor" flag.
func (c *Config) AddEditorFlag(f *pflag.FlagSet) {
	f.StringVar(&c.Editor, EditorFlag, c.Editor, "command-line editor snippet")
}

// AddStartFlags adds all flags related to the "start" subcommand.
func (c *Config) AddStartFlags(f *pflag.FlagSet) {
	c.AddAddrFlag(f)
	c.AddTmpDirFlag(f)
	c.AddEditorFlag(f)
}

// ValidateAddrFlag validates the "addr" flag.
func (c *Config) ValidateAddrFlag() error {
	if c.Addr == "" {
		return fmt.Errorf("%w: flag %q is not informed",
			ErrInvalidConfig, AddrFlag)
	}
	return nil
}

// ValidateEditorFlag validates the "editor" flag.
func (c *Config) ValidateEditorFlag() error {
	if c.Editor == "" {
		return fmt.Errorf("%w: flag %q is not informed",
			ErrInvalidConfig, EditorFlag)
	}
	return nil
}

// ValidateTmpDirFlag validates "tmp-dir" flag.
func (c *Config) ValidateTmpDirFlag() error {
	if c.TmpDir == "" {
		return fmt.Errorf("%w: flag %q is not informed",
			ErrInvalidConfig, TmpDirFlag)
	}
	stat, err := os.Stat(c.TmpDir)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidConfig, err)
	}
	if !stat.IsDir() {
		return fmt.Errorf("%w: %q is not found", ErrInvalidConfig, c.TmpDir)
	}
	return nil
}

// ValidateStartFlags validates all flags employed on "start" subcommand.
func (c *Config) ValidateStartFlags() error {
	var err error
	if err = c.ValidateAddrFlag(); err != nil {
		return err
	}
	if err = c.ValidateEditorFlag(); err != nil {
		return err
	}
	return c.ValidateTmpDirFlag()
}

// LoggerWith decorates logger with the flags informed, where empty flag values
// are skipped.
func (c *Config) LoggerWith(logger *slog.Logger, flags ...string) *slog.Logger {
	m := map[string]string{
		AddrFlag:   c.Addr,
		EditorFlag: c.Editor,
		TmpDirFlag: c.TmpDir,
	}
	for _, k := range flags {
		v, ok := m[k]
		if !ok || v == "" {
			continue
		}
		logger = logger.With(k, v)
	}
	return logger
}

// NewConfig instantiate a new Config with default values.
func NewConfig() *Config {
	defaultLogLevel := slog.LevelDebug
	return &Config{
		LogLevel: &defaultLogLevel,
		Addr:     "127.0.0.1:8928",
		TmpDir:   os.Getenv("TMPDIR"),
		Editor:   os.Getenv("EDITOR"),
	}
}
