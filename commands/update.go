package commands

import (
	"gocker-project/git"
	"gocker-project/yaml"
)

type UpdateCmd struct {
	Services []string `arg:"positional"`
}

func (command *UpdateCmd) Run(config yaml.Config) {
	git := git.NewFromConfig(config)

	if len(command.Services) > 0 {
		filteredRepos := git.FilterRepos(command.Services)
		git.CloneSpecifiedRepos(filteredRepos)
	} else {
		git.CloneAllRepos()
	}
}
