/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     music-grabber
 * @Date:        2024-01-13 00:26
 * @Description:
 */

package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/chromedp/chromedp"
	"github.com/leafney/music-grabber/model"
	"github.com/leafney/music-grabber/pkg/cdper"
	"github.com/leafney/music-grabber/pkg/vars"
	"github.com/leafney/rose"
	"image/color"
	"log"
)

var (
	BrowserUrl = ""
	Version    = "v0.1.0"
)

func main() {

	urlCh := make(chan model.UrlLink, 100)

	a := app.New()
	w := a.NewWindow(fmt.Sprintf("Music Grabber %v", Version))

	//resDataList := binding.BindStringList(&[]string{})

	resList := binding.NewUntypedList()

	go func() {
		for u := range urlCh {
			resList.Prepend(model.UrlLink{
				Url: u.Url,
				Ctx: u.Ctx,
			})
		}
	}()

	searchBtn := widget.NewButton("Open", func() {

		if rose.StrIsEmpty(BrowserUrl) {
			log.Println("未选择")
			return
		}

		go cdper.StartBrowser(BrowserUrl, urlCh)

		//resDataList.Append("hello")

		//resList.Prepend(model.UrlLink{
		//	Url: "1111111",
		//	Ctx: context.Background(),
		//})

	})

	//searchBox := lay2.NewResponsiveLayout(
	//	lay2.Responsive(searchInput, 0.5, 0.75),
	//	lay2.Responsive(searchBtn, 0.2, 0.25),
	//)

	//closeBtn := widget.NewButton("Close", func() {
	//
	//})

	//container.NewVSplit()

	searchLabel := widget.NewLabel("Search:")
	searchInput := widget.NewEntry()
	searchInput.SetPlaceHolder("Please enter the search song title or artist")
	//searchInput.Resize(searchInput.MinSize())
	searchInput.OnChanged = func(s string) {
		if len(s) > 0 {
			searchBtn.SetText("Search")
		} else {
			searchBtn.SetText("Open")
		}
	}

	// 禁用输入框
	searchInput.Disable()

	searchBox := container.NewBorder(
		nil, nil, nil, container.NewHBox(searchBtn), searchInput,
	)

	//:= container.NewHBox(
	//	//widget.NewRichText(),
	//	searchInput,
	//
	//)

	webLabel := widget.NewLabel("Website:")
	//webBox := container.NewHBox(
	//	widget.NewCheck("fangpi", func(v bool) {
	//
	//	}),
	//	widget.NewCheck("NetEaseCloud music", func(v bool) {
	//
	//	}),
	//	widget.NewCheck("QQ music", func(v bool) {
	//
	//	}),
	//)

	radio := widget.NewRadioGroup([]string{"NetEaseCloud", "QQ", "TongZhong", "GeQuBao", "FangPi"}, func(s string) {
		switch s {
		case "NetEaseCloud":
			BrowserUrl = vars.MusicWebHomeNetEaseCloud
		case "QQ":
			BrowserUrl = vars.MusicWebHomeQQ
		case "TongZhong":
			BrowserUrl = vars.MusicWebHomeTonZhon
		case "GeQuBao":
			BrowserUrl = vars.MusicWebHomeGeQuBao
		case "FangPi":
			BrowserUrl = vars.MusicWebHomeFangPi
		default:
			BrowserUrl = vars.MusicWebHomeNetEaseCloud
		}
	})
	radio.SetSelected("NetEaseCloud")
	radio.Horizontal = true

	webBox := container.NewVBox(radio)

	emptyBox := canvas.NewRectangle(color.Transparent)
	emptyBox.Resize(fyne.NewSize(100, 10))

	resultClear := widget.NewButton("Clear Result", func() {
		resList.Set(nil)
	})
	//resultClear.Disable()

	resultLabel := container.NewBorder(nil, nil, nil,
		resultClear, widget.NewLabel("Result:"))

	// test 1
	//resultList := widget.NewListWithData(resDataList,
	//	func() fyne.CanvasObject {
	//		return widget.NewLabel("template")
	//	},
	//	func(item binding.DataItem, obj fyne.CanvasObject) {
	//		obj.(*widget.Label).Bind(item.(binding.String))
	//	},
	//)

	// test 2
	resultBox := widget.NewListWithData(resList,
		func() fyne.CanvasObject {
			return container.NewBorder(
				nil, nil, nil,
				widget.NewButton("Listen in Browser", nil),
				widget.NewLabel(""),
			)
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			//obj.(*widget.Label).Bind(item.(binding.String))

			//ctr, _ := obj.(*fyne.Container)
			//l := ctr.Objects[0].(*widget.Label)
			//b := ctr.Objects[1].(*widget.Button)
			//dui, _ := item.(binding.Untyped).Get()
			//u := dui.(model.UrlLink)
			//l.SetText(u.Url)
			//b.OnTapped = func() {
			//	log.Println("open 点击了", u.Url)
			//	////u.Ctx.
			//	//if err := chromedp.Run(u.Ctx, chromedp.Navigate(u.Url)); err != nil {
			//	//	log.Printf("open error [%v]", err)
			//	//}
			//}

			dui, _ := item.(binding.Untyped).Get()
			u := dui.(model.UrlLink)

			obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(u.Url)
			obj.(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
				log.Println("button click")
				newTabCtx, _ := chromedp.NewContext(u.Ctx)
				if err := chromedp.Run(newTabCtx, chromedp.Navigate(u.Url)); err != nil {
					log.Printf("open error [%v]", err)
				}
			}
		},
	)

	// test 3
	//resultList := widget.NewListWithData(resList,
	//	func() fyne.CanvasObject {
	//		return widget.NewHyperlink("", nil)
	//	},
	//	func(item binding.DataItem, obj fyne.CanvasObject) {
	//		//obj.(*widget.Label).Bind(item.(binding.String))
	//		hyl, _ := obj.(*widget.Hyperlink)
	//		dui := item.(binding.String)
	//		u := dui.(string)
	//
	//		hyl.SetText("")
	//		//hyl.SetURL()
	//
	//	},
	//)

	boxList := container.NewVBox(
		webLabel,
		webBox,
		widget.NewSeparator(),
		searchLabel,
		searchBox,
		emptyBox,
		widget.NewSeparator(),
		resultLabel,
	)

	w.SetContent(container.NewBorder(boxList, nil, nil, nil, resultBox))
	w.Resize(fyne.NewSize(540, 600))
	w.SetFixedSize(true)
	w.ShowAndRun()
}

func showDialog(w fyne.Window) {

	content := container.NewVBox()

	dialog.NewCustom("Tips", "Close", content, w)
}
