package main

import (
	"fmt"

	_ "github.com/CodyGuo/godaemon"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigName(".ups-config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/var/ups-poweroff")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	switch viper.GetString("mode") {

	case "service":
		start_service()
	case "client":
		start_client()
	}
}
