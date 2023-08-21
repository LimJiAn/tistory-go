# tistory-go
[![Go](https://img.shields.io/badge/go-1.19-blue.svg?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/dl/)
[![Chromedp](https://img.shields.io/badge/chromedp-0.9.2-red.svg?style=for-the-badge&logo=go&logoColor=white)](https://pkg.go.dev/github.com/chromedp/chromedp)

> #### tistory-go 는 티스토리 블로그(tistory blog) 자동화를 위한 Go 언어 Library 입니다.


## Installation

As a library

```shell
go get github.com/LimJiAn/tistory-go
```
## Usage

Your Go app you can do something like

```go
package main

import (
    "log"
    "os"

    "github.com/LimJiAn/tistory-go"
)

func main() {
    clientId := "your-client-id"
    clientSecret := "your-client-secret"
    blogURL := "your-blog-url"

    // Create Tistory
    tistory, err := tistory.NewTistory(blogURL, clientId, clientSecret)
    if err != nil {
    	log.Fatal(err)
    }

    // Get AuthorizationCode
    blogId := "your-blog-id"
    blogPassword := "your-blog-password"
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

    // Get Post List
    _, err = tistory.GetPostList(1)
    if err != nil {
        log.Fatal(err)
    }

    // Get Post
    _, err = tistory.GetPost(1)
    if err != nil {
        log.Fatal(err)
    }

    // Write Post
    _, err = tistory.WritePost(
	    map[string]interface{}{"title": "title", "content": "content", "visibility": "3"})
    if err != nil {
        log.Fatal(err)
    }

    // Modify Post
    _, err = tistory.ModifyPost(
        map[string]interface{}{"postId": "1", "title": "title", "content": "content", "visibility": "3"})
    if err != nil {
        log.Fatal(err)
    }
}
```

## Reference
#### [Tistory App Register](https://www.tistory.com/guide/api/manage/register)
#### [Tistory Open API](https://tistory.github.io/document-tistory-apis/)
##
> #### API 사용 중 status 403 , error_message 이 블로그는 내부 정책으로 OPEN API 사용할 수 없습니다.
> #### -> 스팸성 게시물 작성이 증가하여 이용이 제한될 수 있습니다.
