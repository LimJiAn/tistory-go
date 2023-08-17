package main

import (
	"log"
	"os"

	"github.com/LimJiAn/tistory-go"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	clientId := os.Getenv("CLIENT_ID")
	secretKey := os.Getenv("CLIENT_SECRET")
	blogURL := "https://jiaaan90.tistory.com"
	tistory := tistory.NewTistory(blogURL, clientId, secretKey)

	// Get AuthorizationCode
	_, err := tistory.GetAuthorizationCode()
	if err != nil {
		log.Fatal(err)
	}

	// Get AccessToken
	_, err = tistory.GetAccessToken()
	if err != nil {
		log.Fatal(err)
	}
}
