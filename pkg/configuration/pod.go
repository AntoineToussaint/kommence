package configuration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AntoineToussaint/jarvis/pkg/output"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Pod struct {
	ID          string
	Shortcut    string
	Description string
	Service     string
	Namespace   string
	Container   string
	LocalPort   int `yaml:"localPort"`
	PodPort     int `yaml:"podPort"`
}

func NewPod(f string) (*Pod, error) {
	data, err := os.ReadFile(f)
	if err != nil {
		return nil, errors.Wrapf(err, "can't load file: %v", f)
	}
	var cfg Pod
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "can't unmarshal pod configuration")
	}
	cfg.ID = strings.Replace(f, ".yaml", "", 1)
	cfg.ID = strings.Replace(cfg.ID, "jarvis/pods/", "", 1)
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
  name: {{.Service}}
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
		log.Debugf("Pods folder not found in jarvis config\n")
		return &config, nil
	}
	err = filepath.WalkDir(p,
		func(p string, d os.DirEntry, err error) error {
			if !strings.HasSuffix(p, ".yaml") {
				return nil
			}
			c, err := NewPod(p)
			if err != nil {
				log.Errorf("cannot load pod: %v", err)
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
