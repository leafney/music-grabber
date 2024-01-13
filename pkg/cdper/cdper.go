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
	urlCh := make(chan string, 100)

	InterceptResponse(ctx, tIdCh, urlCh)

	//InterceptResponsePlus(ctx, tIdCh)

	//targetCh := chromedp.WaitNewTarget(ctx, func(info *target.Info) bool {
	//	log.Println("打开了新tab")
	//	return info.URL != ""
	//})

	// 启动Chrome浏览器
	err := chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate("https://www.fangpi.net/"),
		//chromedp.Navigate("https://y.qq.com/"),
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
			//log.Printf("[INFO] id [%v] type [%v]", c.Id, c.Type)

			newTabCtx, _ := chromedp.NewContext(c.BaseCtx, chromedp.WithTargetID(c.Id))
			//InterceptResponsePlus(newTabCtx, tIdCh)
			InterceptResponse(newTabCtx, tIdCh, urlCh)
			chromedp.Run(newTabCtx)

			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for u := range urlCh {
			log.Printf("接收到url [%v]", u)
		}
	}()

	select {}
}

func InterceptResponse(ctx context.Context, chT chan model.TargetId, chU chan string) {
	chromedp.ListenTarget(ctx, func(event interface{}) {
		switch ev := event.(type) {
		case *network.EventRequestWillBeSent:
			theUrl := ev.Request.URL
			if ev.Type != network.ResourceTypeMedia {
				break
			}
			if rose.StrAnyContains(theUrl, "mp3", "m4a") {
				//log.Println("2-拦截到音乐文件链接 ", theUrl)
				chU <- theUrl
			}
		case *target.EventTargetCreated:
			c := chromedp.FromContext(ctx)
			go func(theCtx context.Context, baseTarget *chromedp.Target, nowInfo *target.Info) {
				if baseTarget.TargetID == nowInfo.OpenerID {
					chT <- model.TargetId{
						Id:      nowInfo.TargetID,
						BaseCtx: theCtx,
					}
				}
			}(ctx, c.Target, ev.TargetInfo)
		}
	})
}

func InterceptResponsePlus(ctx context.Context, ch chan model.TargetId) {
	eveMap := make(map[string]string)
	chromedp.ListenTarget(ctx, func(event interface{}) {
		switch ev := event.(type) {
		case *network.EventRequestWillBeSent:

			theUrl := ev.Request.URL

			//log.Println("请求url ", theUrl)

			if ev.Type != network.ResourceTypeMedia {
				break
			}

			//log.Println("the media url", theUrl)

			reqID := ev.RequestID
			eveKey := rose.Md5HashStr(reqID.String())

			if rose.StrAnyContains(theUrl, "mp3", "m4a") {
				log.Println("2-拦截到音乐文件链接 ", theUrl)
				eveMap[eveKey] = theUrl

				//

			}

		case *network.EventLoadingFinished:
			reqID := ev.RequestID
			eveKey := rose.Md5HashStr(reqID.String())

			//execCtx := cdp.WithExecutor(theCtx, chromedp.FromContext(theCtx).Target)

			if v, ok := eveMap[eveKey]; ok {
				delete(eveMap, eveKey)

				go func(theCtx context.Context, url string) {
					//	//execCtx := cdp.WithExecutor(theCtx, chromedp.FromContext(theCtx).Target)
					//	if body, err := network.GetResponseBody(ev.RequestID).Do(execCtx); err == nil {
					//		//data := string(body)
					//
					//		//fmt.Println("----------得到内容↓↓↓↓↓↓↓↓↓---------------")
					//		//fmt.Println("body", len(data), len(body))
					//		//fmt.Println("----------得到内容↑↑↑↑↑↑↑↑↑---------------")
					//		//res <- body
					//
					//		//dec, data, _ := minimp3.DecodeFull(body)
					//		//seconds := int64((len(data) - dec.SampleRate) * 8 / (dec.Kbps * 1000))
					//
					//		//dec, _ := mp3.NewDecoder(bytes.NewReader(body))
					//
					//		//streamer, format, err := mp3.Decode(io.NopCloser(bytes.NewBuffer(body)))
					//		//if err != nil {
					//		//	log.Println("error ", err)
					//		//}
					//		//
					//		//length := format.SampleRate.D(streamer.Len())
					//		//log.Println("second", length.Seconds())
					//		//log.Printf("url [%v] size [%v] seconds [%v]", url, utils.FormatFileSize(body), length.Round(time.Second))
					//	}
					//
				}(ctx, v)

			}

			//size := ev.EncodedDataLength
			//log.Println("size ", size)

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
						Type:    1,
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
					Type:    2,
					BaseCtx: ctx,
				}
			}()
		}
	})
}
