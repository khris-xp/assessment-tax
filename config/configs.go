package config

import (
	"os"

	"log"

	"github.com/joho/godotenv"
)

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func EnvPort() string {
	loadEnv()
	return os.Getenv("PORT")
}

func EnvHost() string {
	loadEnv()
	return os.Getenv("HOST")
}

func EnvUser() string {
	loadEnv()
	return os.Getenv("USER")
}

func EnvPassword() string {
	loadEnv()
	return os.Getenv("PASSWORD")
}

func EnvDBName() string {
	loadEnv()
	return os.Getenv("DBNAME")
}

func EnvAdminUsername() string {
	loadEnv()
	return os.Getenv("ADMIN_USERNAME")
}

func EnvAdminPassword() string {
	loadEnv()
	return os.Getenv("ADMIN_PASSWORD")
}
