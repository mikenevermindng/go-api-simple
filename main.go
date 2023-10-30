package main

import (
	"goApiByGin/config"
	"goApiByGin/db"
	"goApiByGin/router"
	customvalidator "goApiByGin/validator"
	"log"
)

func main() {
	if err := config.SetupConfig(); err != nil {
		log.Fatalf("config SetupConfig() error: %s", err)
	}
	customvalidator.SetupValidator()
	connectionString := config.DbConfiguration()
	db.InitialDatabase(connectionString)
	defer db.Client().Close()
	r := router.SetupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
