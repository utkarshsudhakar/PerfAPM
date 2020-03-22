package utils

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"

	"github.com/utkarshsudhakar/demo_app/config"
)

func ReadConfig() config.Config {
	var configfile = "properties.ini"
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	var conf config.Config
	if _, err := toml.DecodeFile(configfile, &conf); err != nil {
		log.Fatal(err)
	}
	//log.Print(config.Index)
	return conf
}
