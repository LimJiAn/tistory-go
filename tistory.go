package tistory

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

type TistoryUser struct {
	LoginId  string
	Password string
}

func AuthTistory() error {
	/*
		chrome windows excute

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

	// create context

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	tistoryAuthURL := "https://www.tistory.com/auth/login"
	if os.Getenv("KAKAO_ID") != "" && os.Getenv("KAKAO_PASSWORD") != "" {
		if err := chromedp.Run(ctx,
			chromedp.Navigate(tistoryAuthURL),
			chromedp.Click(`.btn_login.link_kakao_id`), // class="btn_login link_tistory_id"
			chromedp.Sleep(2*time.Second),
			chromedp.SendKeys(`input[name="loginId"]`, os.Getenv("KAKAO_ID")),
			chromedp.SendKeys(`input[name="password"]`, os.Getenv("KAKAO_PASSWORD")),
			chromedp.Sleep(1*time.Second),
			chromedp.Click(`.btn_g.highlight.submit`), // class="btn_g highlight submit"
			chromedp.Sleep(2*time.Second),
			chromedp.Location(&tistoryAuthURL),
		); err != nil {
			return err
		}
	} else if os.Getenv("TISTORY_ID") != "" && os.Getenv("TISTORY_PASSWORD") != "" {
		if err := chromedp.Run(ctx,
			chromedp.Navigate(tistoryAuthURL),
			chromedp.Click(`.btn_login.link_tistory_id`), // class="btn_login link_tistory_id"
			chromedp.Sleep(2*time.Second),
			chromedp.SendKeys(`input[name="loginId"]`, os.Getenv("TISTORY_ID")),
			chromedp.SendKeys(`input[name="password"]`, os.Getenv("TISTORY_PASSWORD")),
			chromedp.Sleep(1*time.Second),
			chromedp.Click(`btn_login`), // class="btn_login"
			chromedp.Sleep(2*time.Second),
		); err != nil {
			return err
		}
	} else {
		return errors.New(
			"Please TISTORY_ID and TISTORY_PASSWORD or KAKAO_ID and KAKAO_PASSWORD in .env file")
	}
	return nil
}
