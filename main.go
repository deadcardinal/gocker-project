package main

import (
	"gocker-project/commands"
	"gocker-project/yaml"
	"log"

	"github.com/alexflint/go-arg"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Output(0, "[system] .env file not found")
	}
}

func main() {
	var args struct {
		File     string                `default:"docker-compose.yml" arg:"-f" help:"config file"`
		Update   *commands.UpdateCmd   `arg:"subcommand:update"`
		Checkout *commands.CheckoutCmd `arg:"subcommand:checkout"`
		Command  *commands.CommandCmd  `arg:"subcommand:command"`
		Shell    *commands.ShellCmd    `arg:"subcommand:shell"`
		Run      *commands.RunCmd      `arg:"subcommand:run"`
	}
	arg.MustParse(&args)

	config := yaml.GetConfigFromFile(args.File)

	if args.Update != nil {
		args.Update.Run(config)
	}

	if args.Checkout != nil {
		args.Checkout.Run(config)
	}

	if args.Command != nil {
		args.Command.Run(config)
	}

	if args.Shell != nil {
		args.Shell.Run(config)
	}

	if args.Run != nil {
		args.Run.Run(config)
	}
}
