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
	clientSecret := os.Getenv("CLIENT_SECRET")
	blogURL := "https://jiaaan90.tistory.com"
	tistory := tistory.NewTistory(blogURL, clientId, clientSecret)

	blogId := os.Getenv("KAKAO_ID")
	blogPassword := os.Getenv("KAKAO_PASSWORD")
	_, err := tistory.GetAuthorizationCode(blogId, blogPassword)
	if err != nil {
		log.Fatal(err)
	}

	// Get AccessToken
	_, err = tistory.GetAccessToken()
	if err != nil {
		log.Fatal(err)
	}
}
