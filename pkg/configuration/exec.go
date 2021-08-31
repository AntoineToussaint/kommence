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
	Name        string
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
	if cfg.Cmd == ""{
		return nil, nil
	}
	if cfg.Description == "" {
		cfg.Description = "No description available"
	}
	return &cfg, nil
}

func (e Executable) ToString() string {
	return output.FromTemplate(`- {{.Name}}
  Command: {{.Cmd}}
  Description: {{.Description}}
`, e)
}

type Executables struct {
	Commands map[string]*Executable
}

func NewExecutableConfiguration(p string) (*Executables, error) {
	e := Executables{Commands: make(map[string]*Executable)}
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
			e.Commands[c.Name] = c
			shortcut := c.Shortcut
			if shortcut == "" {
				return nil
			}
			if _, ok := e.Commands[shortcut]; ok{
				return fmt.Errorf("shortcut %v duplicated in executables", shortcut)
			}
			e.Commands[shortcut] = c
			return nil
		})
	if err != nil {
		return nil, fmt.Errorf("couldn't find executables at: %v", p)
	}
	return &e, nil
}

func (c *Executables) Get(x string) (*Executable, bool) {
	exec, ok := c.Commands[x]
	return exec, ok
}