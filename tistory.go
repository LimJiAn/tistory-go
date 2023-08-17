package tistory

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type Tistory struct {
	BlogURL            string
	ClientId           string
	SecretKey          string
	LoginURL           string
	AccessToken        string
	AuthenticationURL  string
	RedirectAuthURL    string
	AuthenticationCode string
}

func NewTistory(blogURL, clientId, secretKey string) *Tistory {
	return &Tistory{
		BlogURL:   blogURL,
		ClientId:  clientId,
		SecretKey: secretKey,
		LoginURL: fmt.Sprintf(
			"https://www.tistory.com/auth/login"),
		AuthenticationURL: fmt.Sprintf(
			"https://www.tistory.com/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code", clientId, blogURL),
	}
}

// Login & Get AuthorizationCode
// return authorizationCode, error
func (t *Tistory) GetAuthorizationCode() (string, error) {
	// Excute chrome
	/*
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", false),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-dev-shm-usage", true),
		)
		allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()
		ctx, cancel := chromedp.NewContext(allocCtx)
		defer cancel()
	*/

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Navigate to tistory login page
	if err := chromedp.Run(ctx,
		chromedp.Navigate(t.LoginURL),
	); err != nil {
		return "", err
	}

	// Login
	if os.Getenv("KAKAO_ID") != "" && os.Getenv("KAKAO_PASSWORD") != "" {
		if err := chromedp.Run(ctx,
			// xpath
			chromedp.Click(`//*[@id="cMain"]/div/div/div/a[2]`), // class="btn_login link_tistory_id"
			chromedp.Sleep(2*time.Second),
			chromedp.SendKeys(`//*[@id="loginId--1"]`, os.Getenv("KAKAO_ID")),
			chromedp.SendKeys(`//*[@id="password--2"]`, os.Getenv("KAKAO_PASSWORD")),
			chromedp.Sleep(1*time.Second),
			chromedp.Click(`//*[@id="mainContent"]/div/div/form/div[4]/button[1]`), // class="btn_g highlight submit"
			chromedp.Sleep(2*time.Second),
		); err != nil {
			return "", err
		}
	} else if os.Getenv("TISTORY_ID") != "" && os.Getenv("TISTORY_PASSWORD") != "" {
		if err := chromedp.Run(ctx,
			chromedp.Click(`//*[@id="cMain"]/div/div/div/a[3]`), // class="btn_login link_tistory_id"
			chromedp.Sleep(2*time.Second),
			chromedp.SendKeys(`//*[@id="loginId"]`, os.Getenv("TISTORY_ID")),
			chromedp.SendKeys(`//*[@id="loginPw"]`, os.Getenv("TISTORY_PASSWORD")),
			chromedp.Sleep(1*time.Second),
			chromedp.Click(`//*[@id="authForm"]/fieldset/button`), // class="btn_login"
			chromedp.Sleep(2*time.Second),
		); err != nil {
			return "", err
		}
	} else {
		return "", errors.New(
			"Please KAKAO_ID and KAKAO_PASSWORD or TISTORY_ID and TISTORY_PASSWORD in .env file")
	}

	// Get AuthenticationCode
	if err := chromedp.Run(ctx,
		chromedp.Navigate(t.AuthenticationURL),
		chromedp.Sleep(1*time.Second),
		chromedp.Click(`//*[@id="contents"]/div[4]/button[1]`),
		chromedp.Sleep(1*time.Second),
		chromedp.Location(&t.RedirectAuthURL),
	); err != nil {
		return "", err
	}

	if t.RedirectAuthURL == "" {
		return "", errors.New("Failed to RedirectAuthURL")
	}

	// http://client.redirect.uri?code=authorizationCode&state=someValue
	if len(strings.Split(t.RedirectAuthURL, "code=")) < 2 {
		return "", errors.New("Failed to GetAuthenticationCode")
	}

	t.AuthenticationCode = strings.Split(t.RedirectAuthURL, "code=")[1]
	if len(strings.Split(t.AuthenticationCode, "&state")) < 2 {
		return "", errors.New("Failed to GetAuthenticationCode")
	}

	t.AuthenticationCode = strings.Split(t.AuthenticationCode, "&state")[0]
	return t.AuthenticationCode, nil
}

func (t *Tistory) GetAccessToken() (string, error) {
	params := url.Values{
		"client_id":     {t.ClientId},
		"client_secret": {t.SecretKey},
		"redirect_uri":  {t.BlogURL},
		"code":          {t.AuthenticationCode},
		"grant_type":    {"authorization_code"},
	}

	accessTokenURL := fmt.Sprintf(
		"https://www.tistory.com/oauth/access_token?%s", params.Encode())
	resp, err := http.Get(accessTokenURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	bodyString := string(respBytes)
	if len(strings.Split(bodyString, "=")) < 2 {
		return "", errors.New(
			"Failed to GetAccessToken from bodyString (len(strings.Split(bodyString, \"=\")) < 2)")
	}

	if !strings.HasPrefix(bodyString, "access_token") {
		return "", errors.New(
			"Failed to GetAccessToken from bodyString (strings.Contains(bodyString, \"access_token\")")
	}

	t.AccessToken = strings.Split(bodyString, "=")[1]
	return t.AccessToken, nil
}
