package manipulator

import (
	"gocker-project/yaml"
	"strings"
	"sync"
)

type Manipulator struct {
	Commands map[int]Command
}

func NewFromConfig(config yaml.Config) Manipulator {
	manipulator := Manipulator{
		Commands: map[int]Command{},
	}

	i := 0

	for service, serviceData := range config.Services {
		command := Command{}

		for label, labelValue := range serviceData.Labels {
			if label != "project.git" && label != "project.git.branch" && strings.HasPrefix(label, "project.") {
				firstDot := strings.Index(label, ".") + 1
				command.Name = label[firstDot:len(label)]
				command.Exec = labelValue
				command.ServiceName = service
				command.Path = "./src/" + service
			}
		}

		if (Command{} != command) {
			manipulator.Commands[i] = command
		}

		i++
	}

	return manipulator
}

func (manipulator *Manipulator) RunForSpecifiedServices(commandName string, serviceNames []string) {
	var wg sync.WaitGroup
	// wg.Add(len(serviceNames))

	for _, serviceName := range serviceNames {
		for _, command := range manipulator.Commands {
			if command.Name == commandName && command.ServiceName == serviceName {
				wg.Add(1)
				go func(command Command) {
					command.Run()
					wg.Done()
				}(command)
			}
		}
	}

	wg.Wait()
}
