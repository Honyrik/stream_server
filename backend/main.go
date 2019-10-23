package main

import (
	"os"
	filepath "path"
	"strings"

	logger "github.com/jeanphorn/log4go"
)

func main() {
	fileConfig := os.Getenv("STREAM_SERVER_CONFIG")
	var err error
	if len(strings.TrimSpace(fileConfig)) == 0 {
		dirConfig, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		fileConfig = filepath.Join(dirConfig, "config.toml")
	}
	logger.LoadConfiguration("./logger.json")
	logger.Info("Init config")
	var config Config
	if config, err = initConfig(fileConfig); err != nil {
		panic(err)
	}
	startHTTP(config)
}
