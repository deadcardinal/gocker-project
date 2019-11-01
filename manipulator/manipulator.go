package manipulator

import (
	"fmt"
	"gocker-project/yaml"
	"strings"
)

type Manipulator struct {
	Commands map[string]Command
}

func NewFromConfig(config yaml.Config) Manipulator {
	manipulator := Manipulator{
		Commands: map[string]Command{},
	}

	for service, serviceData := range config.Services {
		command := Command{}

		for label, labelValue := range serviceData.Labels {
			if label != "project.git" && label != "project.git.branch" && strings.HasPrefix(label, "project.") {
				firstDot := strings.Index(label, ".") + 1
				command.Name = label[firstDot:len(label)]
				command.Exec = labelValue
				command.ServiceName = service
				command.Path = "src/" + service
			}
		}

		if (Command{} != command) {
			manipulator.Commands[service] = command
		}
	}

	fmt.Printf("%v\n", manipulator.Commands)

	return manipulator
}

// func (manipulator *Manipulator) run(services []Service, name string) {

// }
