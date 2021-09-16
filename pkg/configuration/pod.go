package configuration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/AntoineToussaint/kommence/pkg/output"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Pod struct {
	ID          string
	Shortcut    string
	Description string
	Name        string
	Namespace   string
	Container   string
	LocalPort   int `yaml:"localPort"`
	PodPort     int `yaml:"podPort"`
}

func NewPod(f string) (*Pod, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, errors.Wrapf(err, "can't load file: %v", f)
	}
	var cfg Pod
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "can't unmarshal flow configuration")
	}
	if cfg.Name == "" {
		return nil, nil
	}
	if cfg.Namespace == "" {
		return nil, nil
	}
	if cfg.Description == "" {
		cfg.Description = "No description available"
	}
	return &cfg, nil
}

func (p Pod) ToString(log *output.Logger) string {
	return output.FromTemplate(log, `- {{.ID}}
  name: {{.Name}}
  namespace: {{.Namespace}}
  {{if .Container}}container: {{.Container}}{{end}}
  port: {{.LocalPort}} -> {{.PodPort}}
  Description: {{.Description}}
`, p)
}

type Pods struct {
	Pods      map[string]*Pod
	Shortcuts map[string]*Pod
}

func NewPodConfiguration(log *output.Logger, p string) (*Pods, error) {
	config := Pods{Pods: make(map[string]*Pod), Shortcuts: make(map[string]*Pod)}
	dir, err := os.Stat(p)
	if err != nil || !dir.IsDir() {
		log.Debugf("Pods folder not found in kommence config")
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
			c, err := NewPod(p)
			if err != nil {
				return err
			}
			config.Pods[c.ID] = c
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
		return nil, fmt.Errorf("couldn't find executables at: %v", p)
	}
	return &config, nil
}

func (c *Pods) Get(x string) (*Pod, bool) {
	exec, ok := c.Pods[x]
	if !ok {
		exec, ok = c.Shortcuts[x]
	}
	return exec, ok
}
