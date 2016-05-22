package utils

import (
	"github.com/BurntSushi/toml"
)

var (
	log = SetUpLogger("config")
)

type Config struct {
	DBUrl            string
	DBName           string
	DBCollectionName string
	AESKey           string
}

var Conf Config

func SetConfig() {
	_, err := toml.DecodeFile("config.toml", &Conf)
	if err != nil {
		log.Error("Error parsing config:", err)
		panic(err)
	}
	log.Infof("Read config file: mongos url: %s, AESKey: %s, main db name: %s, main collection name: %s\n", Conf.DBUrl, Conf.AESKey, Conf.DBName, Conf.DBCollectionName)
}
