package pkg

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type FlowConfiguration struct {
	Runs []*RunnerConfiguration
	Kubes []*KubeConfiguration

}

func LoadRun(f string, name string) (*RunnerConfiguration, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, errors.Wrapf(err, "can't load file: %v", f)
	}
	var cfg RunnerConfiguration
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "can't unmarshal flow configuration")
	}
	if cfg.Run.Cmd == "" {
		return nil, nil
	}
	// Overwrite the name from the file name if not set
	if cfg.Run.Name == "" {
		cfg.Run.Name = strings.Split(name, ".")[0]
	}
	return &cfg, nil
}


func LoadKube(f string, name string) (*KubeConfiguration, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, errors.Wrapf(err, "can't load file: %v", f)
	}
	var cfg KubeConfiguration
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "can't unmarshal flow configuration")
	}
	if cfg.Kube.Deployment == "" {
		return nil, nil
	}
	// Overwrite the name from the file name if not set
	if cfg.Kube.Name == "" {
		cfg.Kube.Name = strings.Split(name, ".")[0]
	}
	return &cfg, nil
}

func LoadFlowConfiguration(kommenceDir string) (*FlowConfiguration, error) {
	// We load all the files in the kommence directory
	config := FlowConfiguration{}
	err := filepath.Walk(kommenceDir,
		func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			run, err := LoadRun(p, info.Name())
			if err != nil {
				log.Println(err)
				return err
			}
			if run != nil {
				config.Runs = append(config.Runs, run)
			}
			kube, err := LoadKube(p, info.Name())
			if err != nil {
				log.Println(err)
				return err
			}
			if kube != nil {
				config.Kubes = append(config.Kubes, kube)
			}
			return nil
		})
	if err != nil {
		return nil, errors.Wrap(err, "can't load kommence configuration")
	}
	// load kubernetes client
	if len(config.Kubes) > 0 {
		LoadKubeClient()
	}
	return &config, nil
}


