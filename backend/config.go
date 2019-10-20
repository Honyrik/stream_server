package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/naoina/toml"
)

type Config struct {
	VideoDir    string
	SiteDir     string
	Listen      string
	AcceptsFile string
	VideoTmp    string
}

func initConfig(file string) (Config, error) {
	var config Config
	if _, err := os.Stat(file); os.IsNotExist(err) {
		os.Create(file)
		return config, fmt.Errorf("Not found " + file)
	}
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return config, err
	}

	toml.Unmarshal(dat, &config)
	fileold := file + ".bak"
	f, err := os.Create(fileold)
	if err != nil {
		return config, err
	}
	encoder := toml.NewEncoder(f)
	encoder.Encode(config)
	defer f.Close()
	return config, nil
}
