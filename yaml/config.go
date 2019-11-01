package yaml

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Services map[string]Service
}

// type Service struct {
// 	Labels   map[string]string `yaml:"labels"`
// 	Repo     git.Repo
// 	Commands Command
// }

// type Command struct {
// 	ServiceName string
// 	Name        string
// 	Path        string
// 	Exec        string
// }

// type Manipulator struct {
// 	Commands map[string]Command
// }

// func runCommand(config Config) {
// 	var manipulator = parseManipulator(config)
// 	manipulator.Commands = nil
// }

func GetConfigFromYaml(file string) Config {
	data, err := ioutil.ReadFile(file)

	log.Output(0, "[system] use file "+file)
	if err != nil {
		log.Fatal(err)
	}

	config := Config{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

// func parseManipulator(config Config) Manipulator {
// 	manipulator := Manipulator{
// 		Commands: map[string]Command{},
// 	}

// 	for service, serviceData := range config.Services {
// 		command := Command{}

// 		for label, labelValue := range serviceData.Labels {
// 			if label != "project.git" && label != "project.git.branch" && strings.HasPrefix(label, "project.") {
// 				firstDot := strings.Index(label, ".") + 1
// 				command.Name = label[firstDot:len(label)]
// 				command.Exec = labelValue
// 				command.ServiceName = service
// 				command.Path = "src/" + service
// 			}
// 		}

// 		if (Command{} != command) {
// 			manipulator.Commands[service] = command
// 		}
// 	}

// 	fmt.Printf("%v\n", manipulator.Commands)

// 	return manipulator
// }

// func (manipulator *Manipulator) run(services []Service, name string) {

// }
