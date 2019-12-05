package main

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
)

func main() {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://xxxholic.top/login.html"),
		chromedp.WaitVisible(`.userName`),
		chromedp.SendKeys(`.userName`, "sa"),
		chromedp.WaitVisible(`.passWord`),
		chromedp.SendKeys(`.passWord`, "sa"),
		chromedp.Click(`#login`),
	); err != nil {
		fmt.Println(err)
	}
	fmt.Println("登陆成功")
}
