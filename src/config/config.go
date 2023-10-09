package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

var (
	DbStringConnection = ""
	Port               = 0
	SecretKey          []byte
)

func Load() {
	var err error

	if os.Getenv("ENVIRONMENT") != "PROD" {
		if err = godotenv.Load(); err != nil {
			log.Fatal(err)
		}
	}

	Port, err = strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		Port = 8000
	}

	DbStringConnection = fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
	)

	SecretKey = []byte(os.Getenv("SECRET_KEY"))
}
