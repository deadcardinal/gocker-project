package commands

import (
	"gocker-project/manipulator"
	"gocker-project/yaml"
)

type CommandCmd struct {
	Command string `arg:"positional" arg:"required"`
}

func (command *CommandCmd) Run(config yaml.Config) {
	var manipulator = manipulator.NewFromConfig(config)
	manipulator.Commands = nil
}
