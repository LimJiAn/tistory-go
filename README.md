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

#### Authorization ([인증 및 권한](https://tistory.github.io/document-tistory-apis/auth/authorization_code.html))
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
}
```
#### BlogInfo ([블로그 정보](https://tistory.github.io/document-tistory-apis/apis/v1/blog/list.html))
```go
    // Blog Info
    info, err := tistory.GetBlogInfo()
    if err != nil {
        log.Fatal(err)
    }
```

#### PostList ([글 목록](https://tistory.github.io/document-tistory-apis/apis/v1/post/list.html))
```go
    // Post List
    res, err := tistory.GetPostList(1)
    if err != nil {
        log.Fatal(err)
    }
```

#### ReadPost ([글 읽기](https://tistory.github.io/document-tistory-apis/apis/v1/post/read.html))
```go
    // Read Post
    res, err := tistory.GetPost(1)
    if err != nil {
        log.Fatal(err)
    }
```

#### WritePost ([글 작성](https://tistory.github.io/document-tistory-apis/apis/v1/post/write.html))
```go
    // Write Post
    res, err := tistory.WritePost(
        map[string]interface{}{"title": "title", "content": "content", "visibility": "3"})
    if err != nil {
        log.Fatal(err)
    }
```

#### ModifyPost ([글 수정](https://tistory.github.io/document-tistory-apis/apis/v1/post/modify.html))
```go
    // Modify Post
    res, err := tistory.ModifyPost(
        map[string]interface{}{"postId": "1", "title": "title", "content": "content", "visibility": "3"})
    if err != nil {
        log.Fatal(err)
    }
```

#### AttchFile ([파일 첨부](https://tistory.github.io/document-tistory-apis/apis/v1/post/attach.html))
```go
    // Attach File (only image)
    fileName := "/UserFilepath/test.png"
    res, err := tistory.AttachPost(fileName)
    if err != nil {
        log.Fatal(err)
    }
```

#### CategoryList ([카테고리 목록](https://tistory.github.io/document-tistory-apis/apis/v1/category/list.html))
```go
    // Category List
    res, err := tistory.CategoryList()
    if err != nil {
        log.Fatal(err)
    }
```

#### RecentComment ([최근 댓글 목록](https://tistory.github.io/document-tistory-apis/apis/v1/comment/recent.html))
```go
    // Recent Comment List
    res, err := tistory.GetRecentCommentList(1, 1)
    if err != nil {
        log.Fatal(err)
    }
```

#### CommentList ([게시글 댓글 목록](https://tistory.github.io/document-tistory-apis/apis/v1/comment/list.html))
```go
    // Comment List
    res, err := tistory.GetCommentList(1)
    if err != nil {
        log.Fatal(err)
    }
```

#### WriteComment ([댓글 작성](https://tistory.github.io/document-tistory-apis/apis/v1/comment/write.html))
```go
    // Write Comment
    res, err := tistory.WriteComment(
        map[string]interface{}{"postId": "1", "content": "comment"})
    if err != nil {
        log.Fatal(err)
    }
```

#### ModifyComment ([댓글 수정](https://tistory.github.io/document-tistory-apis/apis/v1/comment/modify.html))
```go
    // Modify Comment
    info, err := tistory.ModifyComment(
        map[string]interface{}{"postId": "1", "commentId": "1", "content": "comment"})
    if err != nil {
        log.Fatal(err)
    }
```

#### DeleteComment ([댓글 삭제](https://tistory.github.io/document-tistory-apis/apis/v1/comment/delete.html))
```go
    // Delete Comment
    info, err := tistory.DeleteComment(
        map[string]interface{}{"postId": "1", "commentId": "1"})
    if err != nil {
        log.Fatal(err)
    }
```

## Reference
#### [Tistory App Register](https://www.tistory.com/guide/api/manage/register)
#### [Tistory Open API](https://tistory.github.io/document-tistory-apis/)
##
> #### API 사용 중 status 403 , error_message 이 블로그는 내부 정책으로 OPEN API 사용할 수 없습니다.
> #### -> 스팸성 게시물 작성이 증가하여 이용이 제한될 수 있습니다.
