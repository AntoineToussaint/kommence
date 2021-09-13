package configuration

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/antoinetoussaint/kommence/pkg/output"
	"github.com/pkg/errors"
)

type Flow struct {
	ID          string
	Shortcut    string
	Description string
	Executables []string
	Pods []string
}

func NewFlow(f string) (*Flow, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, errors.Wrapf(err, "can't load file: %v", f)
	}
	var cfg Flow
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "can't unmarshal flow configuration")
	}
	if cfg.ID == "" {
		return nil, fmt.Errorf("ID required")
	}
	if cfg.Description == "" {
		cfg.Description = "No description available"
	}
	return &cfg, nil
}

func (f *Flow) ToString(log *output.Logger) string {
	return output.FromTemplate(log, `- {{.ID}}
  Description: {{.Description}}
`, f)
}

type Flows struct {
	Flows  map[string]*Flow
	Shortcuts map[string]*Flow
}

func NewFlowConfiguration(p string) (*Flows, error) {
	config := Flows{Flows: make(map[string]*Flow), Shortcuts: make(map[string]*Flow)}
	err := filepath.Walk(p,
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			c, err := NewFlow(p)
			if err != nil {
				return err
			}
			config.Flows[c.ID] = c
			shortcut := c.Shortcut
			if shortcut == "" {
				return nil
			}
			if _, ok := config.Shortcuts[shortcut]; ok {
				return fmt.Errorf("shortcut %v duplicated in flows", shortcut)
			}
			config.Shortcuts[shortcut] = c
			return nil
		})
	if err != nil {
		return nil, fmt.Errorf("error loading configurations %v: %v", p, err)
	}
	return &config, nil
}

func (c *Flows) Get(x string) (*Flow, bool) {
	flow, ok := c.Flows[x]
	if !ok {
		flow, ok = c.Shortcuts[x]
	}
	return flow, ok
}

func (c *Flows) GetExecutables(x string) []string {
	flow, ok := c.Flows[x]
	if !ok {
		flow, ok = c.Shortcuts[x]
	}
	if flow == nil {
		return nil
	}
	return flow.Executables
}

func (c *Flows) GetPods(x string) []string {
	flow, ok := c.Flows[x]
	if !ok {
		flow, ok = c.Shortcuts[x]
	}
	if flow == nil {
		return nil
	}
	return flow.Pods
}
