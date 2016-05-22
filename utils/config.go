package utils

import (
	"github.com/BurntSushi/toml"
	"fmt"
)
var (
	log = SetUpLogger("config")
)
type Config struct {
	DBUrl string
	DBName string
	DBCollectionName string
	AESKey string
}

var Conf Config

func SetConfig() {
	_, err := toml.DecodeFile("config.toml", &Conf); if err != nil {
		fmt.Println("Error parsing config")
		panic(err)
	}
	fmt.Println("Read config file: mongos url: " + Conf.DBUrl + ", AESKey: " + Conf.AESKey + ", main db name: " + Conf.DBName + ", main collection name: " + Conf.DBCollectionName)
}