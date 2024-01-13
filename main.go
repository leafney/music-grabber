/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     music-grabber
 * @Date:        2024-01-13 00:26
 * @Description:
 */

package main

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/leafney/music-grabber/model"
	"log"
)

func main() {
	a := app.New()
	w := a.NewWindow("Music Grabber")

	//resDataList := binding.BindStringList(&[]string{})

	resList := binding.NewUntypedList()

	searchBtn := widget.NewButton("Open", func() {
		//cdper.StartBrowser()

		//resDataList.Append("hello")

		resList.Prepend(model.UrlLink{
			Url: "1111111",
			Ctx: context.Background(),
		})
	})

	//searchBox := lay2.NewResponsiveLayout(
	//	lay2.Responsive(searchInput, 0.5, 0.75),
	//	lay2.Responsive(searchBtn, 0.2, 0.25),
	//)

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

	searchBox := container.NewBorder(
		nil, nil, nil, searchBtn, searchInput,
	)

	//:= container.NewHBox(
	//	//widget.NewRichText(),
	//	searchInput,
	//
	//)

	webLabel := widget.NewLabel("Music Website:")
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

	radio := widget.NewRadioGroup([]string{"NetEaseCloud music", "QQ music", "FangPi"}, func(s string) {
		log.Println("choose", s)
	})
	radio.SetSelected("NetEaseCloud music")
	radio.Horizontal = true

	webBox2 := container.NewHBox(radio)

	//hello := widget.NewLabel("Hello Fyne!")
	//(container.NewVBox(
	//
	//	widget.NewButton("Browser", func() {
	//		//hello.SetText("Welcome!")
	//		cdper.StartBrowser()
	//	}),
	//))

	resultClear := widget.NewButton("Clear Result", func() {
		//resDataList.Set([]string{})
		resList.Set([]interface{}{})
	})
	//resultClear.Disable()

	resultLabel := container.NewBorder(nil, nil, nil,
		resultClear, widget.NewLabel("Result:"))

	//resultList := widget.NewListWithData(resDataList,
	//	func() fyne.CanvasObject {
	//		return widget.NewLabel("template")
	//	},
	//	func(item binding.DataItem, obj fyne.CanvasObject) {
	//		obj.(*widget.Label).Bind(item.(binding.String))
	//	},
	//)

	resultList := widget.NewListWithData(resList,
		func() fyne.CanvasObject {
			return container.NewBorder(
				nil, nil, nil,
				widget.NewButton("Open in Browser", func() {

				}),
				widget.NewLabel(""),
			)
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			//obj.(*widget.Label).Bind(item.(binding.String))
			ctr, _ := obj.(*fyne.Container)
			l := ctr.Objects[0].(*widget.Label)
			b := ctr.Objects[1].(*widget.Button)
			dui, _ := item.(binding.Untyped).Get()
			u := dui.(model.UrlLink)
			l.SetText(u.Url)
			b.OnTapped = func() {
				//u.Ctx.
			}
		},
	)

	boxList := container.NewVBox(
		webLabel,
		//webBox,
		webBox2,
		searchLabel,
		searchBox,

		resultLabel,
		//resultBox,
	)

	w.SetContent(container.NewBorder(boxList, nil, nil, nil, resultList))
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}

func ToLaunch(c context.Context) func() {
	return func() {
		fmt.Println(c)
	}
}
