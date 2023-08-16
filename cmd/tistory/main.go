package main

import (
	"log"

	"github.com/LimJiAn/tistory"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// Login to tistory and get authorization_code
	err := tistory.AuthTistory()
	if err != nil {
		log.Fatal(err)
	}
}
