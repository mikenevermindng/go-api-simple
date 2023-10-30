package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type DatabaseConfiguration struct {
	Driver   string
	Dbname   string
	Username string
	Password string
	Host     string
	Port     string
	LogMode  bool
}

func DbConfiguration() string {
	DBName := viper.GetString("DB_NAME")
	DBUser := viper.GetString("DB_USER")
	DBPassword := viper.GetString("DB_PASSWORD")
	DBHost := viper.GetString("DB_HOST")
	DBPort := viper.GetString("DB_PORT")

	DBDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		DBHost, DBUser, DBPassword, DBName, DBPort,
	)

	return DBDSN
}
