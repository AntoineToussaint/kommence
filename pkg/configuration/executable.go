package configuration

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"

	"github.com/AntoineToussaint/kommence/pkg/output"
	"github.com/pkg/errors"
)

// Executable configuration.
type Executable struct {
	ID          string
	Shortcut    string
	Description string
	Cmd         string
	Env         map[string]string
	Watch       []string
	StdErr      string `yaml:"std_err"`
}

const (
	Ignore  = "ignore"
	AsError = "error"
	AsLog   = "log"
)

// NewExecutable attempts to load a configuration.
func NewExecutable(f string) (*Executable, error) {
	data, err := os.ReadFile(f)
	if err != nil {
		return nil, errors.Wrapf(err, "can't load file: %v", f)
	}
	var cfg Executable
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "can't unmarshal flow configuration")
	}
	if cfg.ID == "" {
		return nil, fmt.Errorf("ID required")
	}
	if cfg.Cmd == "" {
		return nil, fmt.Errorf("command required")
	}
	if cfg.StdErr == "" {
		cfg.StdErr = Ignore
	}
	if cfg.Description == "" {
		cfg.Description = "No description available"
	}
	return &cfg, nil
}

// ToString convert to string.
func (e *Executable) ToString(log *output.Logger) string {
	return output.FromTemplate(log, `- {{.ID}}
  command: {{.Cmd}}
  Description: {{.Description}}
`, e)
}

// Executables aggregate Executable configurations.
type Executables struct {
	Commands  map[string]*Executable
	Shortcuts map[string]*Executable
}

// NewExecutableConfiguration loads Executables configuration.
func NewExecutableConfiguration(log *output.Logger, p string) (*Executables, error) {
	config := Executables{Commands: make(map[string]*Executable), Shortcuts: make(map[string]*Executable)}
	dir, err := os.Stat(p)
	if err != nil || !dir.IsDir() {
		log.Debugf("Executables folder not found in kommence config\n")
		return &config, nil
	}
	err = filepath.Walk(p,
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			c, err := NewExecutable(p)
			if err != nil {
				return err
			}
			config.Commands[c.ID] = c
			shortcut := c.Shortcut
			if shortcut == "" {
				return nil
			}
			if _, ok := config.Shortcuts[shortcut]; ok {
				return fmt.Errorf("shortcut %v duplicated in executables", shortcut)
			}
			config.Shortcuts[shortcut] = c
			return nil
		})
	if err != nil {
		return nil, fmt.Errorf("error loading configurations %v: %v", p, err)
	}
	return &config, nil
}

// Get an Executable by ID or shortcut.
func (c *Executables) Get(x string) (*Executable, bool) {
	exec, ok := c.Commands[x]
	if !ok {
		exec, ok = c.Shortcuts[x]
	}
	return exec, ok
}
