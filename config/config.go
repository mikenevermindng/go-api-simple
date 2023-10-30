package config

import (
	"github.com/spf13/viper"
	"log"
)

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
}

// SetupConfig configuration
func SetupConfig() error {
	var configuration *Configuration

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error to reading config file, %s", err)
		return err
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("error to decode, %v", err)
		return err
	}

	return nil
}
