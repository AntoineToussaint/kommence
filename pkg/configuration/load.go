package configuration

import (
	"fmt"
	"path"
	"strings"

	"github.com/antoinetoussaint/kommence/pkg/output"
	"github.com/fatih/color"
)

type Configuration struct {
	Execs *Executables
	Pods  *Pods
	Flows *Flows
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

	// Flows configurations
	flows, err := NewFlowConfiguration(path.Join(p, "/flows"))
	if err != nil {
		return nil, fmt.Errorf("can't load flows configurations: %v", err)
	}
	cfg.Flows = flows
	logger.Debugf("loaded %v flow configurations\n", len(flows.Flows))

	return &cfg, nil
}

func (c *Configuration) ListExecutables() []string {
	if len(c.Execs.Commands) == 0 {
		return nil
	}
	var execs []string
	for _, exec := range c.Execs.Commands {
		s := exec.ID
		if shortcut := exec.Shortcut; shortcut != "" {
			s = fmt.Sprintf("%v|%v", s, shortcut)
		}
		execs = append(execs, s)
	}
	return execs
}


func (c *Configuration) ValidExecutables(execs []string) (bool, string) {
	var unknowns []string
	for _, exec := range execs {
		if _, ok := c.Execs.Get(exec); !ok {
			unknowns = append(unknowns, exec)
		}
	}
	if len(unknowns) > 0 {
		return false, "Unknown executables " + strings.Join(unknowns, ", ")
	}
	return true, ""
}

func (c *Configuration) ListPods() []string {
	if len(c.Pods.Pods) == 0 {
		return nil
	}
	var pods []string
	for _, pod := range c.Pods.Pods {
		s := pod.ID
		if shortcut := pod.Shortcut; shortcut != "" {
			s = fmt.Sprintf("%v|%v", s, shortcut)
		}
		pods = append(pods, s)
	}
	return pods
}


func (c *Configuration) ValidPods(pods []string) (bool, string) {
	var unknowns []string
	for _, pod := range pods {
		if _, ok := c.Pods.Get(pod); !ok {
			unknowns = append(unknowns, pod)
		}
	}
	if len(unknowns) > 0 {
		return false, "Unknown pods " + strings.Join(unknowns, ", ")
	}
	return true, ""
}


func (c *Configuration) ListFlows() []string {
	if len(c.Flows.Flows) == 0 {
		return nil
	}
	var flows []string
	for _, flow := range c.Flows.Flows {
		s := flow.ID
		if shortcut := flow.Shortcut; shortcut != "" {
			s = fmt.Sprintf("%v|%v", s, shortcut)
		}
		flows = append(flows, s)
	}
	return flows
}


func (c *Configuration) ValidFlows(flows []string) (bool, string) {
	var unknowns []string
	for _, flow := range flows {
		if _, ok := c.Flows.Get(flow); !ok {
			unknowns = append(unknowns, flow)
		}
	}
	if len(unknowns) > 0 {
		return false, "Unknown flows " + strings.Join(unknowns, ", ")
	}
	return true, ""
}



func (c *Configuration) Print(logger *output.Logger) {
	logger.Printf("Configured with %v executables:\n", len(c.Execs.Commands), color.Bold)
	for _, exec := range c.Execs.Commands {
		logger.Printf(exec.ToString(logger))
	}
	logger.Printf("Configured with %v pods:\n", len(c.Pods.Pods), color.Bold)
	for _, pod := range c.Pods.Pods {
		logger.Printf(pod.ToString(logger))
	}
	logger.Printf("Configured with %v flows:\n", len(c.Flows.Flows), color.Bold)
	for _, flow := range c.Flows.Flows {
		logger.Printf(flow.ToString(logger))
	}
}
