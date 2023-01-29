package cmd

import (
	"log"
	"sooprox/types"
	"sooprox/utils"
)

func Init(config types.Config, isCli bool) {
	configChan := make(chan types.Config)

	if !isCli {
		log.Printf("Watching config file %s", config.ConfigFile)
		go utils.WatchConfig(configChan, config.ConfigFile)
	}

	Listen(config, configChan)

}
