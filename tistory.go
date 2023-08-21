// Package tistory implements a Tistory API Client.
//
// Tistory API Reference: https://tistory.github.io/document-tistory-apis/apis/

package tistory

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
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
	BlogName           string
	ClientId           string
	ClientSecret       string
	AccessToken        string
	AuthenticationURL  string
	RedirectAuthURL    string
	AuthenticationCode string
}

func NewTistory(blogURL, clientId, clientSecret string) (*Tistory, error) {
	if blogURL == "" || clientId == "" || clientSecret == "" {
		return nil, errors.New("blogURL or clientId or clientSecret is empty")
	}

	if len(strings.Split(blogURL, "//")) < 2 {
		return nil, errors.New("blogURL is invalid")
	}

	if !strings.HasPrefix(blogURL, "https://") {
		return nil, errors.New("blogURL is invalid")
	}

	return &Tistory{
		BlogURL:      blogURL,
		BlogName:     strings.Split(strings.Split(blogURL, "//")[1], ".")[0],
		ClientId:     clientId,
		ClientSecret: clientSecret,
		AuthenticationURL: fmt.Sprintf(
			"https://www.tistory.com/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code", clientId, blogURL),
	}, nil
}

/*
Login & Get AuthorizationCode
https://tistory.github.io/document-tistory-apis/auth/authorization_code.html
*/
func (t *Tistory) GetAuthorizationCode(id, password string) (string, error) {
	if id == "" || password == "" {
		return "", errors.New("id or password is empty")
	}

	/*
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", false),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-dev-shm-usage", true),
		)

		allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()
	*/

	// Create chrome instance
	ctx, cancel := chromedp.NewContext(context.Background())
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
		chromedp.Sleep(1*time.Second),
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

/*
GetAccessToken
https://tistory.github.io/document-tistory-apis/auth/authorization_code.html
*/
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

/*
GetBlogInfo 블로그 정보
access_token: 발급받은 access_token
output: 출력방식
https://tistory.github.io/document-tistory-apis/apis/v1/blog/list.html
*/
func (t *Tistory) GetBlogInfo() (map[string]interface{}, error) {
	t.AccessToken = "ce196a9e476dd617519f2074286aa5c9_6a5d731ffc6022608f75393a7ea87cf4"
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

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf("Failed to GetBlogInfo (resp.StatusCode: %d)", resp.StatusCode))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

/*
GetPostList 글 목록
blogName: Blog Name (필수)
page: 불러올 페이지 번호
https://tistory.github.io/document-tistory-apis/apis/v1/post/list.html
*/
func (t *Tistory) GetPostList(pageNumber int) (map[string]interface{}, error) {
	params := url.Values{
		"access_token": {t.AccessToken},
		"blogName":     {t.BlogName},
		"page":         {fmt.Sprintf("%d", pageNumber)},
	}

	postListURL := fmt.Sprintf(
		"https://www.tistory.com/apis/post/list?%s", params.Encode())
	resp, err := http.Get(postListURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf("Failed to GetPostList (resp.StatusCode: %d)", resp.StatusCode))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

/*
GetPost 글 읽기
blogName: Blog Name (필수)
postId: 글 번호 (필수)
https://tistory.github.io/document-tistory-apis/apis/v1/post/read.html
*/
func (t *Tistory) GetPost(postId int) (map[string]interface{}, error) {
	params := url.Values{
		"access_token": {t.AccessToken},
		"blogName":     {t.BlogName},
		"postId":       {fmt.Sprintf("%d", postId)}}

	postURL := fmt.Sprintf(
		"https://www.tistory.com/apis/post/read?%s", params.Encode())

	resp, err := http.Get(postURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf("Failed to GetPost (resp.StatusCode: %d)", resp.StatusCode))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

/*
WritePost 글 작성
blogName: Blog Name (필수)
title: 글 제목 (필수)
content: 글 내용
visibility: 발행상태 (0: 비공개 - 기본값, 1: 보호, 3: 발행)
category: 카테고리 아이디 (기본값: 0)
published: 발행시간 (TIMESTAMP 이며 미래의 시간을 넣을 경우 예약. 기본값: 현재시간)
slogan: 문자 주소
tag: 태그 (',' 로 구분)
acceptComment: 댓글 허용 (0, 1 - 기본값)
password: 보호글 비밀번호
https://tistory.github.io/document-tistory-apis/apis/v1/post/write.html
*/
func (t *Tistory) WritePost(option map[string]interface{}) (map[string]interface{}, error) {
	params := url.Values{
		"access_token": {t.AccessToken},
		"output":       {"json"},
		"blogName":     {t.BlogName},
	}

	for key, value := range option {
		params.Add(key, fmt.Sprintf("%v", value))
	}

	writePostURL := "https://www.tistory.com/apis/post/write?"
	resp, err := http.PostForm(writePostURL, params)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf("Failed to WritePost (resp.StatusCode: %d)", resp.StatusCode))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

/*
ModifyPost 글 수정
blogName: Blog Name (필수)
postId: 글 번호 (필수)
title: 글 제목 (필수)
content: 글 내용
visibility: 발행상태 (0: 비공개 - 기본값, 1: 보호, 3: 발행)
category: 카테고리 아이디 (기본값: 0)
published: 발행시간 (TIMESTAMP 이며 미래의 시간을 넣을 경우 예약. 기본값: 현재시간)
slogan: 문자 주소
tag: 태그 (',' 로 구분)
acceptComment: 댓글 허용 (0, 1 - 기본값)
password: 보호글 비밀번호
https://tistory.github.io/document-tistory-apis/apis/v1/post/modify.html
*/
func (t *Tistory) ModifyPost(option map[string]interface{}) (map[string]interface{}, error) {
	params := url.Values{
		"access_token": {t.AccessToken},
		"output":       {"json"},
		"blogName":     {t.BlogName},
	}

	for key, value := range option {
		params.Add(key, fmt.Sprintf("%v", value))
	}

	modifyPostURL := "https://www.tistory.com/apis/post/modify?"
	resp, err := http.PostForm(modifyPostURL, params)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf("Failed to ModifyPost (resp.StatusCode: %d)", resp.StatusCode))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

/*
AttachPost 파일 첨부
blogName: Blog Name
uploadedfile: 업로드할 파일 (multipart/form-data)
https://tistory.github.io/document-tistory-apis/apis/v1/post/attach.html
*/
func (t *Tistory) AttachPost(uploadedfile *multipart.FileHeader) (map[string]interface{}, error) {
	params := url.Values{
		"access_token": {t.AccessToken},
		"blogName":     {t.BlogName},
	}

	content, err := uploadedfile.Open()
	defer content.Close()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to AttachPost")
	}

	attachPostURL := fmt.Sprintf(
		"https://www.tistory.com/apis/post/attach?%s", params.Encode())

	resp, err := http.Post(attachPostURL, "multipart/form-data", content)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf("Failed to AttachPost (resp.StatusCode: %d)", resp.StatusCode))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

/*
CategoryList 카테고리 목록
id: 카테고리 ID
name: 카테고리 이름
parent: 부모 카테고리 ID
label: 부모 카테고리를 포함한 전체 이름 ('/'로 구분)
entries: 카테고리내 글 수
*/
func (t *Tistory) CategoryList() (map[string]interface{}, error) {
	params := url.Values{
		"access_token": {t.AccessToken},
		"blogName":     {t.BlogName},
		"output":       {"json"},
	}

	categoryListURL := fmt.Sprintf(
		"https://www.tistory.com/apis/category/list?%s", params.Encode())

	resp, err := http.Get(categoryListURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf("Failed to CategoryList (resp.StatusCode: %d)", resp.StatusCode))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

/*
GetComment 최신 댓글 목록 가져오기
blogName: Blog Name
page: 가져올 페이지 (기본값: 1)
count: 페이지당 댓글 수 (기본값: 10, 최대값: 10)
https://tistory.github.io/document-tistory-apis/apis/v1/comment/recent.html
*/
func (t *Tistory) GetNewCommentList(page, count int) (map[string]interface{}, error) {
	params := url.Values{
		"access_token": {t.AccessToken},
		"blogName":     {t.BlogName},
		"output":       {"json"},
		"page":         {fmt.Sprintf("%d", page)},
		"count":        {fmt.Sprintf("%d", count)},
	}

	getNewCommentListURL := fmt.Sprintf(
		"https://www.tistory.com/apis/comment/newest?%s", params.Encode())

	resp, err := http.Get(getNewCommentListURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf("Failed to GetNewCommentList (resp.StatusCode: %d)", resp.StatusCode))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

/*
GetCommentList 댓글 목록 가져오기
blogName: Blog Name
postId: 글 번호 (ID)
https://tistory.github.io/document-tistory-apis/apis/v1/comment/list.html
*/
func (t *Tistory) GetCommentList(postId int) (map[string]interface{}, error) {
	params := url.Values{
		"access_token": {t.AccessToken},
		"blogName":     {t.BlogName},
		"output":       {"json"},
		"postId":       {fmt.Sprintf("%d", postId)},
	}

	getCommentListURL := fmt.Sprintf(
		"https://www.tistory.com/apis/comment/list?%s", params.Encode())

	resp, err := http.Get(getCommentListURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf("Failed to GetCommentList (resp.StatusCode: %d)", resp.StatusCode))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

/*
WriteComment 댓글 작성
blogName: Blog Name (필수)
postId: 글 ID (필수)
parentId: 부모 댓글 ID (대댓글인 경우 사용)
content: 댓글 내용
secret: 비밀 댓글 여부 (1: 비밀댓글, 0: 공개댓글 - 기본 값)
https://tistory.github.io/document-tistory-apis/apis/v1/comment/write.html
*/
func (t *Tistory) WriteComment(option map[string]interface{}) (map[string]interface{}, error) {
	params := url.Values{
		"access_token": {t.AccessToken},
		"blogName":     {t.BlogName},
	}

	for key, value := range option {
		params.Add(key, fmt.Sprintf("%v", value))
	}

	writeCommentURL := "https://www.tistory.com/apis/comment/write?"
	resp, err := http.PostForm(writeCommentURL, params)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf("Failed to WriteComment (resp.StatusCode: %d)", resp.StatusCode))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

/*
ModifyComment 댓글 수정
blogName: Blog Name (필수)
postId: 글 ID (필수)
commentId: 댓글 ID (필수)
parentId: 부모 댓글 ID (대댓글인 경우 사용)
content: 댓글 내용
secret: 비밀 댓글 여부 (1: 비밀댓글, 0: 공개댓글 - 기본 값)
https://tistory.github.io/document-tistory-apis/apis/v1/comment/modify.html
*/
func (t *Tistory) ModifyComment(option map[string]interface{}) (map[string]interface{}, error) {
	params := url.Values{
		"access_token": {t.AccessToken},
		"blogName":     {t.BlogName},
	}

	for key, value := range option {
		params.Add(key, fmt.Sprintf("%v", value))
	}

	modifyCommentURL := "https://www.tistory.com/apis/comment/modify?"
	resp, err := http.PostForm(modifyCommentURL, params)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf("Failed to ModifyComment (resp.StatusCode: %d)", resp.StatusCode))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

/*
DeleteComment 댓글 삭제
blogName: Blog Name (필수)
postId: 글 ID (필수)
commentId: 댓글 ID (필수)
https://tistory.github.io/document-tistory-apis/apis/v1/comment/delete.html
*/
func (t *Tistory) DeleteComment(option map[string]interface{}) (map[string]interface{}, error) {
	params := url.Values{
		"access_token": {t.AccessToken},
		"output":       {"json"},
		"blogName":     {t.BlogName},
	}

	for key, value := range option {
		params.Add(key, fmt.Sprintf("%v", value))
	}

	deleteCommentURL := "https://www.tistory.com/apis/comment/delete?"
	resp, err := http.PostForm(deleteCommentURL, params)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(
			fmt.Sprintf("Failed to DeleteComment (resp.StatusCode: %d)", resp.StatusCode))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
