package utils

import (
	"log"
	"os"
	"sooprox/types"

	"gopkg.in/yaml.v3"
)

func ReadConfig(config ...string) types.Config {
	_config := "config.yaml"
	if len(config) > 0 {
		_config = config[0]
	}
	f, err := os.ReadFile(_config)
	if err != nil {
		// log.Fatal(err)
		log.Printf("Error reading config file: %s", err)
		return types.Config{}
	}

	// Create an empty Car to be are target of unmarshalling
	var c types.Config

	// Unmarshal our input YAML file into empty Car (var c)
	if err := yaml.Unmarshal(f, &c); err != nil {
		log.Printf("Error unmarshalling config file: %s", err)
		return types.Config{}
	}

	return c
}
