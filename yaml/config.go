package yaml

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Services map[string]Service
}

func GetConfigFromFile(file string) Config {
	data, err := ioutil.ReadFile(file)

	log.Output(0, "[system] use file "+file)

	if err != nil {
		log.Fatal(err)
	}

	return GetConfigFromBytes(data)
}

func GetConfigFromBytes(data []byte) Config {
	config := Config{}
	err := yaml.Unmarshal(data, &config)

	if err != nil {
		log.Fatal(err)
	}

	return config
}
