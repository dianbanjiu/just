package config

import (
	"github.com/spf13/viper"
	"log"
)

var CFG = viper.New()

func init() {
	CFG.SetConfigName("config")
	CFG.SetConfigType("yaml")
	CFG.AddConfigPath(".")
	err := CFG.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}
