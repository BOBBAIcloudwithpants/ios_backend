package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func init() {
	fmt.Println("viper")
	viper.SetConfigName("configure")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
