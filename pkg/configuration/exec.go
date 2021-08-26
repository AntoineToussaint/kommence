package configuration

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Executable struct {
	Cmd string
	Name string
	Shortcut string
	Watch []string
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
	return &cfg, nil
}

type Executables struct {
	execs []*Executable
}

func NewExecutableConfiguration(p string) (*Executables, error) {
	e := Executables{}
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
			e.execs = append(e.execs, c)
			return nil
		})
	if err != nil {
		return nil, fmt.Errorf("couldn't find executables at: %v", p)
	}
	return &e, nil
}