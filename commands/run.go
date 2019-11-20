package commands

import (
	"gocker-project/docker"
	"gocker-project/yaml"
)

type RunCmd struct {
	Services []string `arg:"positional" arg:"required"`
}

func (command *RunCmd) Run(config yaml.Config) {
	if len(command.Services) > 0 {
		filteredServices := docker.FilterServices(config, command.Services)
		docker.Run(filteredServices)
	}
}
