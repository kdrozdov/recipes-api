package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func LoadConfig(path string) {
	viper.SetConfigName("env")
	viper.SetConfigType("json")
	viper.AddConfigPath(path)

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
