/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     music-grabber
 * @Date:        2024-01-13 00:42
 * @Description:
 */

package cdper

import (
	"context"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"github.com/leafney/music-grabber/model"
	"github.com/leafney/music-grabber/pkg/vars"
	"github.com/leafney/rose"
	"log"
	"time"
)

func StartBrowser() {
	opts := []chromedp.ExecAllocatorOption{

		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.ExecPath(vars.MacOSChromePath),
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36"),
	}

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	// 初始化浏览器上下文
	ctx, cancel = chromedp.NewContext(ctx) //chromedp.WithDebugf(log.Printf),

	defer cancel()

	//chromedp.ListenTarget(ctx, InterceptResponse())

	tIdCh := make(chan model.TargetId, 100)

	InterceptResponsePlus(ctx, tIdCh)

	//targetCh := chromedp.WaitNewTarget(ctx, func(info *target.Info) bool {
	//	log.Println("打开了新tab")
	//	return info.URL != ""
	//})

	// 启动Chrome浏览器
	err := chromedp.Run(ctx,
		network.Enable(),
		//chromedp.Navigate("https://www.fangpi.net/"),
		chromedp.Navigate("https://y.qq.com/"),
	)
	if err != nil {
		log.Fatal(err)
	}

	//time.Sleep(10 * time.Second)

	//for {
	//	select {
	//	case tid := <-targetCh:
	//		go func() {
	//			//newTagCtx, cancel := chromedp.NewContext(ctx, chromedp.WithTargetID(tid))
	//			//defer cancel()
	//			log.Println("获取到了新的tab", tid)
	//
	//			//chromedp.Run(newTagCtx)
	//		}()
	//	}
	//}
	//
	//targets, err := chromedp.Targets(ctx)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//for _, t := range targets {
	//	newTabCtx, _ := chromedp.NewContext(ctx, chromedp.WithTargetID(t.TargetID))
	//	//defer newTabCancel()
	//	chromedp.ListenTarget(newTabCtx, InterceptResponse())
	//	log.Println("设置listen ", t.URL)
	//	chromedp.Run(newTabCtx)
	//}

	go func() {
		for c := range tIdCh {
			log.Printf("[INFO] id [%v] type [%v]", c.Id, c.Type)

			newTabCtx, _ := chromedp.NewContext(c.BaseCtx, chromedp.WithTargetID(c.Id))
			InterceptResponsePlus(newTabCtx, tIdCh)
			chromedp.Run(newTabCtx)

			time.Sleep(1 * time.Second)
		}
	}()

	select {}
}

func InterceptResponse() func(event interface{}) {
	return func(event interface{}) {
		if ev, ok := event.(*network.EventResponseReceived); ok {
			if ev.Type != network.ResourceTypeMedia {
				return
			}

			respUrl := ev.Response.URL
			log.Println("response url", respUrl)

			if rose.StrAnyContains(respUrl, "mp3", "m4a") {

				log.Println("1-拦截到音乐文件链接 ", respUrl)

			}
		}
	}
}

func InterceptResponsePlus(ctx context.Context, ch chan model.TargetId) {
	//eveMap := make(map[string]string)
	chromedp.ListenTarget(ctx, func(event interface{}) {
		switch ev := event.(type) {
		case *network.EventRequestWillBeSent:

			theUrl := ev.Request.URL

			log.Println("请求url ", theUrl)

			if ev.Type != network.ResourceTypeMedia {
				break
			}

			log.Println("the media url", theUrl)

			if rose.StrAnyContains(theUrl, "mp3", "m4a") {
				log.Println("2-拦截到音乐文件链接 ", theUrl)
			}

			//reqID := ev.RequestID
			//eveKey := rose.Md5HashStr(reqID.String())
			//eveMap[eveKey]=

		case *network.EventLoadingFinished:
		//reqID := ev.RequestID
		//eveKey := rose.Md5HashStr(reqID.String())

		//execCtx := cdp.WithExecutor(theCtx, chromedp.FromContext(theCtx).Target)

		//go func() {
		//	c := chromedp.FromContext(ctx)
		//	ctx := cdp.WithExecutor(ctx, c.Target)
		//	if body, err := network.GetResponseBody(ev.RequestID).Do(ctx); err == nil {
		//		data := string(body)
		//
		//		fmt.Println("----------得到内容↓↓↓↓↓↓↓↓↓---------------")
		//		fmt.Println(data)
		//		fmt.Println("----------得到内容↑↑↑↑↑↑↑↑↑---------------")
		//		//res <- body
		//	}
		//}()

		case *page.EventWindowOpen:
			c := chromedp.FromContext(ctx)

			//theCtx := cdp.WithExecutor(ctx, chromedp.FromContext(ctx).Target)
			//chromedp.Targets(theCtx)
			log.Println("打开了新窗口", c.Target.TargetID)

		case *target.EventTargetCreated:
			c := chromedp.FromContext(ctx)

			tid := ev.TargetInfo.TargetID
			log.Println("EventTargetCreated", tid)

			if ev.TargetInfo.OpenerID == c.Target.TargetID {
				log.Println("EventTargetCreated equal")
				go func() {
					ch <- model.TargetId{
						Id:      tid,
						Type:    "created",
						BaseCtx: ctx,
					}
				}()
			}

		case *target.EventTargetDestroyed:
			tid := ev.TargetID
			log.Println("EventTargetDestroyed", tid)
			go func() {
				ch <- model.TargetId{
					Id:      tid,
					Type:    "destroyed",
					BaseCtx: ctx,
				}
			}()
		}
	})
}
