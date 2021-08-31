package configuration

import (
	"fmt"
	"github.com/antoinetoussaint/kommence/pkg/output"
	"github.com/fatih/color"
	"path"
)

type Configuration struct {
	Execs *Executables

}

func Load(logger *output.Logger, p string) (*Configuration, error) {
	cfg := Configuration{}
	execs, err := NewExecutableConfiguration(path.Join(p, "/exec"))
	if err != nil {
		return nil, fmt.Errorf("can't load executable configurations: %v", err)
	}
	cfg.Execs = execs
	logger.Debugf("loaded %v executable configurations", len(execs.Commands))
	return &cfg, nil
}


func (c *Configuration) Print(logger *output.Logger) {
	logger.Printf("Configured with %v executables:\n", len(c.Execs.Commands), color.Bold)
	for _, exec := range c.Execs.Commands {
		logger.Printf(exec.ToString())
	}
}



