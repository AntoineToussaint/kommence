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

// Flow is a combination of Executables and Pods
type Flow struct {
	ID          string
	Shortcut    string
	Description string
	Executables []string
	Pods        []string
}

// NewFlow attempts to load a configuration.
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

// ToString converts to string.
func (f *Flow) ToString(log *output.Logger) string {
	return output.FromTemplate(log, `- {{.ID}}
  Description: {{.Description}}
`, f)
}

// Flows aggregate Flow configurations.
type Flows struct {
	Flows     map[string]*Flow
	Shortcuts map[string]*Flow
}

// NewFlowConfiguration loads Flows configuration.
func NewFlowConfiguration(log *output.Logger, p string) (*Flows, error) {
	config := Flows{Flows: make(map[string]*Flow), Shortcuts: make(map[string]*Flow)}
	dir, err := os.Stat(p)
	if err != nil || !dir.IsDir() {
		log.Debugf("Flows folder not found in kommence config")
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

// Get a Flow by ID or shortcut.
func (c *Flows) Get(x string) (*Flow, bool) {
	flow, ok := c.Flows[x]
	if !ok {
		flow, ok = c.Shortcuts[x]
	}
	return flow, ok
}

// GetExecutables get Executables from a Flow by ID or shortcut.
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

// GetPods get Pods from a Flow by ID or shortcut.
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
