package configuration

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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
	Delay       string
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
	if cfg.Cmd == "" {
		return nil, fmt.Errorf("command required")
	}
	if cfg.StdErr == "" {
		cfg.StdErr = Ignore
	}
	if cfg.Description == "" {
		cfg.Description = "No description available"
	}
	cfg.ID = strings.Replace(f, ".yaml", "", 1)
	cfg.ID = strings.Replace(cfg.ID, "kommence/executables/", "", 1)
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
	log.Debugf("walking %s", p)
	err = filepath.WalkDir(p, func(s string, d fs.DirEntry, err error) error {
		if !strings.HasSuffix(s, ".yaml") {
			return nil
		}
		c, err := NewExecutable(s)
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
