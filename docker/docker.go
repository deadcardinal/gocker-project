package docker

import (
	"gocker-project/executor"
	"gocker-project/yaml"
	"strings"
)

func Run(services map[string]yaml.Service) {
	args := make([]string, 0)
	args = append(args, "up", "-d")

	for service, data := range services {
		args = append(args, service)
		for label, labelValue := range data.Labels {
			if strings.HasPrefix(label, "project.depend") {
				args = append(args, labelValue)
			}
		}
	}

	if len(args) > 2 {
		executor.Run("docker-compose", ".", "docker-compose", args...)
	}
}

func FilterServices(config yaml.Config, services []string) map[string]yaml.Service {
	configServices := make(map[string]yaml.Service, 0)

	for _, service := range services {
		if _, ok := config.Services[service]; ok {
			configServices[service] = config.Services[service]
		}
	}

	return configServices
}
