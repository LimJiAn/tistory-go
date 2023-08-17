// Package tistory implements a Tistory API Client.
//
// Tistory API Reference: https://tistory.github.io/document-tistory-apis/apis/

package tistory

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
)

const (
	loginURL = `https://www.tistory.com/auth/login`

	loginKakaoButtonXPath  = `//*[@id="cMain"]/div/div/div/a[2]`
	loginKaKaoIdXPath      = `//*[@id="loginId--1"]`
	loginKakaoPwXPath      = `//*[@id="password--2"]`
	submitKaKaoButtonXPath = `//*[@id="mainContent"]/div/div/form/div[4]/button[1]`

	loginTistoryButtonXPath  = `//*[@id="cMain"]/div/div/div/a[3]`
	loginTistoryIdXPath      = `//*[@id="loginId"]`
	loginTistoryPwXPath      = `//*[@id="loginPw"]`
	submitTistoryButtonXPath = `//*[@id="authForm"]/fieldset/button`

	authButtonXPath = `//*[@id="contents"]/div[4]/button[1]`

	loginAfterURL = `https://www.tistory.com/`
)

type Tistory struct {
	BlogURL            string
	ClientId           string
	ClientSecret       string
	AccessToken        string
	AuthenticationURL  string
	RedirectAuthURL    string
	AuthenticationCode string
}

func NewTistory(blogURL, clientId, clientSecret string) *Tistory {
	return &Tistory{
		BlogURL:      blogURL,
		ClientId:     clientId,
		ClientSecret: clientSecret,
		AuthenticationURL: fmt.Sprintf(
			"https://www.tistory.com/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code", clientId, blogURL),
	}
}

// Login & Get AuthorizationCode
// return authorizationCode, error
func (t *Tistory) GetAuthorizationCode(id, password string) (string, error) {
	if id == "" || password == "" {
		return "", errors.New("id or password is empty")
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Create chrome instance
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Navigate to tistory login page
	if err := chromedp.Run(ctx,
		chromedp.Navigate(loginURL),
	); err != nil {
		return "", errors.Wrap(err, "Failed to Navigate to tistory login page")
	}

	var confirmURL string
	// KAKAO_ID, KAKAO_PASSWORD
	if err := chromedp.Run(ctx,
		chromedp.Click(loginKakaoButtonXPath), // class="btn_login link_tistory_id"
		chromedp.Sleep(1*time.Second),
		chromedp.SendKeys(loginKaKaoIdXPath, id),
		chromedp.SendKeys(loginKakaoPwXPath, password),
		chromedp.Sleep(2*time.Second),
		chromedp.Click(submitKaKaoButtonXPath), // class="btn_g highlight submit"
		chromedp.Sleep(2*time.Second),
		chromedp.Location(&confirmURL),
	); err != nil {
		return "", errors.Wrap(err, "Failed to Login with KAKAO_ID, KAKAO_PASSWORD")
	}

	if confirmURL != loginAfterURL {
		return "", errors.New("Failed to Login")
	}

	// Get AuthenticationCode
	if err := chromedp.Run(ctx,
		chromedp.Navigate(t.AuthenticationURL),
		chromedp.Sleep(1*time.Second),
		chromedp.Click(authButtonXPath),
		chromedp.Sleep(1*time.Second),
		chromedp.Location(&t.RedirectAuthURL),
	); err != nil {
		return "", errors.Wrap(err, "Failed to GetAuthenticationCode")
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
		"client_secret": {t.ClientSecret},
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

func (t *Tistory) GetBlogInfo() (string, error) {
	params := url.Values{
		"access_token": {t.AccessToken},
		"output":       {"json"},
	}

	blogInfoURL := fmt.Sprintf(
		"https://www.tistory.com/apis/blog/info?%s", params.Encode())
	resp, err := http.Get(blogInfoURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	bodyString := string(respBytes)
	return bodyString, nil
}
