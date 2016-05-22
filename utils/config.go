package utils

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	DBUrl string
	DBName string
	DBCollectionName string
	AESKey string
	DBUsersCollectionName string
}

var Conf Config

func SetConfig() {
	toml.DecodeFile("config.toml", &Conf)
}