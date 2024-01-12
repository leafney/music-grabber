/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     music-grabber
 * @Date:        2024-01-13 00:42
 * @Description:
 */

package cdp

import (
	"context"
	"github.com/chromedp/chromedp"
	"log"
	"time"
)

func StartBrowser() {
	ctx, cancel := chromedp.NewExecAllocator(context.Background(),
		chromedp.ExecPath("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"),
	)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()
	// 启动Chrome浏览器
	err := chromedp.Run(ctx, chromedp.Navigate("https://www.baidu.com"))
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(10 * time.Second)
}
