package commands

import (
	"gocker-project/manipulator"
	"gocker-project/yaml"
)

type CommandCmd struct {
	Command  string   `arg:"positional" arg:"required"`
	Services []string `arg:"positional"`
}

func (command *CommandCmd) Run(config yaml.Config) {
	manipulator := manipulator.NewFromConfig(config)

	if len(command.Services) > 0 {
		manipulator.RunForSpecifiedServices(command.Command, command.Services)
	} else {

		serviceNames := make([]string, 0, len(config.Services))
		for serviceName := range config.Services {
			serviceNames = append(serviceNames, serviceName)
		}

		manipulator.RunForSpecifiedServices(command.Command, serviceNames)
	}

}
