package configuration

import (
	"fmt"
	"github.com/antoinetoussaint/kommence/pkg/output"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Executable struct {
	Cmd         string
	ID          string
	Shortcut    string
	Description string
	Watch       []string
}

func NewExecutable(f string) (*Executable, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, errors.Wrapf(err, "can't load file: %v", f)
	}
	var cfg Executable
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "can't unmarshal flow configuration")
	}
	if cfg.ID == ""{
		return nil, fmt.Errorf("ID required")
	}
	if cfg.Cmd == ""{
		return nil, fmt.Errorf("command required")
	}
	if cfg.Description == "" {
		cfg.Description = "No description available"
	}
	return &cfg, nil
}

func (e Executable) ToString() string {
	return output.FromTemplate(`- {{.ID}}
  command: {{.Cmd}}
  Description: {{.Description}}
`, e)
}

type Executables struct {
	Commands map[string]*Executable
	Shortcuts map[string]*Executable
}

func NewExecutableConfiguration(p string) (*Executables, error) {
	config := Executables{Commands: make(map[string]*Executable), Shortcuts: make(map[string]*Executable)}
	err := filepath.Walk(p,
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
			if _, ok := config.Shortcuts[shortcut]; ok{
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

func (c *Executables) Get(x string) (*Executable, bool) {
	exec, ok := c.Commands[x]
	if !ok {
		exec, ok = c.Shortcuts[x]
	}
	return exec, ok
}