package configuration

import (
	"fmt"
	"path"

	"github.com/antoinetoussaint/kommence/pkg/output"
	"github.com/fatih/color"
)

type Configuration struct {
	Execs *Executables
	Pods  *Pods
}

func Load(logger *output.Logger, p string) (*Configuration, error) {
	cfg := Configuration{}

	// Executable configurations
	execs, err := NewExecutableConfiguration(path.Join(p, "/executables"))
	if err != nil {
		return nil, fmt.Errorf("can't load executable configurations: %v", err)
	}
	cfg.Execs = execs
	logger.Debugf("loaded %v executable configurations\n", len(execs.Commands))

	// Pod configurations
	pods, err := NewPodConfiguration(path.Join(p, "/pods"))
	if err != nil {
		return nil, fmt.Errorf("can't load pods configurations: %v", err)
	}
	cfg.Pods = pods
	logger.Debugf("loaded %v pod configurations\n", len(pods.Pods))

	return &cfg, nil
}

func (c *Configuration) Print(logger *output.Logger) {
	logger.Printf("Configured with %v executables:\n", len(c.Execs.Commands), color.Bold)
	for _, exec := range c.Execs.Commands {
		logger.Printf(exec.ToString(logger))
	}
	logger.Printf("Configured with %v pods:\n", len(c.Pods.Pods), color.Bold)
	for _, pod := range c.Pods.Pods {
		fmt.Println(pod)
		logger.Printf(pod.ToString(logger))
	}
}
