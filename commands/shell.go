package commands

import (
	"bufio"
	"fmt"
	"gocker-project/env"
	"gocker-project/executor"
	"gocker-project/yaml"
	"log"
	"os"
	"strings"
	"sync"
)

type ShellCmd struct {
	Services []string `arg:"positional"`
}

func (command *ShellCmd) Run(config yaml.Config) {
	appConfig := env.GetConfig()
	reader := bufio.NewReader(os.Stdin)
	for {
		if len(command.Services) > 0 {
			choicedServices := strings.Join(command.Services, ",")
			fmt.Print("[" + choicedServices + "] \u25b6 ")

		} else {
			fmt.Print("[all] \u25b6 ")
		}

		input, err := reader.ReadString('\n')
		if err != nil {
			log.Output(0, "[system] "+err.Error())
		}

		if len(command.Services) > 0 {
			var wg sync.WaitGroup
			wg.Add(len(command.Services))

			for _, serviceName := range command.Services {
				go func(serviceName string, input string) {
					path := fmt.Sprintf("./%s/%s", appConfig.SourceDir, serviceName)
					if err = executor.ExecInput(serviceName, path, input); err != nil {
						log.Output(0, "["+serviceName+"] "+err.Error())
					}
					wg.Done()
				}(serviceName, input)
			}

			wg.Wait()
		} else {
			var wg sync.WaitGroup
			wg.Add(len(config.Services))

			for serviceName := range config.Services {
				go func(serviceName string, input string) {
					path := fmt.Sprintf("./%s/%s", appConfig.SourceDir, serviceName)
					if err = executor.ExecInput(serviceName, path, input); err != nil {
						log.Output(0, "["+serviceName+"] "+err.Error())
					}
					wg.Done()
				}(serviceName, input)
			}

			wg.Wait()
		}
	}
}
