package main

import (
	"gocker-project/commands"
	"gocker-project/yaml"

	"github.com/alexflint/go-arg"
)

func main() {
	var args struct {
		File     string                `default:"docker-compose.yml" arg:"-f" help:"config file"`
		Update   *commands.UpdateCmd   `arg:"subcommand:update"`
		Checkout *commands.CheckoutCmd `arg:"subcommand:checkout"`
		Command  *commands.CommandCmd  `arg:"subcommand:command"`
	}
	arg.MustParse(&args)

	config := yaml.GetConfigFromYaml(args.File)

	if args.Update != nil {
		args.Update.Run(config)
	}

	if args.Checkout != nil {
		args.Checkout.Run(config)
	}

	if args.Command != nil {
		args.Command.Run(config)
	}
}
