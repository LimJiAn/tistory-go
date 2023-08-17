package main

import (
	"fmt"
	"log"
	"os"

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
	clientId := os.Getenv("CLIENT_ID")
	secretKey := os.Getenv("SECRET_KEY")
	blogURL := "https://jiaaan90.tistory.com"
	tistory := tistory.NewTistory(blogURL, clientId, secretKey)

	authorizationCode, err := tistory.GetAuthorizationCode()
	if err != nil {
		panic(err)
	}
	fmt.Printf("authorizationCode: %v\n", authorizationCode)
}
