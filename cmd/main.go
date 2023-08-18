package main

import (
	"log"
	"os"

	"github.com/LimJiAn/tistory-go"
)

func main() {
	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	blogURL := "https://jiaaan90.tistory.com"

	// Create Tistory
	tistory, err := tistory.NewTistory(blogURL, clientId, clientSecret)
	if err != nil {
		log.Fatal(err)
	}

	// Get AuthorizationCode
	blogId := os.Getenv("KAKAO_ID")
	blogPassword := os.Getenv("KAKAO_PASSWORD")
	_, err = tistory.GetAuthorizationCode(blogId, blogPassword)
	if err != nil {
		log.Fatal(err)
	}

	// Get AccessToken
	_, err = tistory.GetAccessToken()
	if err != nil {
		log.Fatal(err)
	}

	// Get Blog Info
	_, err = tistory.GetBlogInfo()
	if err != nil {
		log.Fatal(err)
	}
}
