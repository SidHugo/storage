package utils

import (
	"github.com/BurntSushi/toml"
)

var (
	log    = SetUpLogger("config")
	format = "Read config file: mongos url: %s, AESKey: %s, main db name: %s, main collection name: %s, users collection name: %s, auth login: %s, auth password: %s"
)

type Config struct {
	DBUrl                 string
	DBName                string
	DBCollectionName      string
	AESKey                string
	DBUsersCollectionName string
	AuthLogin             string
	AuthPassword          string
}

var Conf Config

func SetDefaultConfig() {
	_, err := toml.DecodeFile("config.toml", &Conf)
	if err != nil {
		log.Error("Error parsing config: ", err)
		panic(err)
	}

	log.Infof(format, Conf.DBUrl, Conf.AESKey, Conf.DBName, Conf.DBCollectionName, Conf.DBUsersCollectionName, Conf.AuthLogin, Conf.AuthPassword)
}

func SetConfig(filepath string) {
	_, err := toml.DecodeFile(filepath, &Conf)
	if err != nil {
		log.Error("Error parsing config: ", err)
		panic(err)
	}

	log.Infof(format, Conf.DBUrl, Conf.AESKey, Conf.DBName, Conf.DBCollectionName, Conf.DBUsersCollectionName, Conf.AuthLogin, Conf.AuthPassword)
}
