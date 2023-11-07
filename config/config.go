package config

import (
	"log"
	"os"
	"path"
	"runtime"

	"github.com/spf13/viper"
)

var CFG = viper.New()

func init() {
	CFG.SetConfigName("config")
	CFG.SetConfigType("yaml")
	CFG.AddConfigPath(".")
	var in string
	switch runtime.GOOS {
	case "linux", "darwin":
		in = path.Join(os.Getenv("HOME"), ".config", "just")
	case "windows":
		in = path.Join(os.Getenv("appdata"), "just")
	}
	CFG.AddConfigPath(in)
	err := CFG.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}
