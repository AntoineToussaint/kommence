package pkg

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type FlowConfiguration struct {
	Run []RunnerConfiguration
}

func LoadFlowConfiguration(cfg *viper.Viper) (*FlowConfiguration, error){
	var config FlowConfiguration
	err := cfg.Unmarshal(&config)
	if err != nil {
		return nil, errors.Wrap(err, "can't unmarshal flow configuration")
	}
	return &config, nil
}


