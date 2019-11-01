package commands

import (
	"gocker-project/git"
	"gocker-project/yaml"
)

type CheckoutCmd struct {
	Branch   string   `arg:"positional"`
	Services []string `arg:"positional"`
}

func (command *CheckoutCmd) Run(config yaml.Config) {
	git := git.NewFromConfig(config)

	if len(command.Services) > 0 {
		filteredRepos := git.FilterRepos(command.Services)
		git.CheckoutSpecifiedRepos(filteredRepos, command.Branch)
	} else {
		git.CheckoutAllRepos(command.Branch)
	}
}
