/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     music-grabber
 * @Date:        2024-01-13 00:26
 * @Description:
 */

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	lay2 "fyne.io/x/fyne/layout"
	"log"
)

func main() {
	a := app.New()
	w := a.NewWindow("Music Grabber")

	searchLabel := widget.NewLabel("Search:")
	searchInput := widget.NewEntry()
	searchInput.SetPlaceHolder("Please enter the search song title or artist")
	//searchInput.Resize(searchInput.MinSize())

	data := binding.BindStringList(&[]string{"item1", "item2", "item3"})

	searchBtn := widget.NewButton("Search", func() {
		//cdper.StartBrowser()

		data.Append("hello")
	})

	searchBox := lay2.NewResponsiveLayout(
		lay2.Responsive(searchInput, 0.5, 0.75),
		lay2.Responsive(searchBtn, 0.2, 0.25),
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

	webBox2 := container.NewHBox(radio)

	//hello := widget.NewLabel("Hello Fyne!")
	//(container.NewVBox(
	//
	//	widget.NewButton("Browser", func() {
	//		//hello.SetText("Welcome!")
	//		cdper.StartBrowser()
	//	}),
	//))

	resultLabel := widget.NewLabel("Result:")

	resultList := widget.NewListWithData(data,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			obj.(*widget.Label).Bind(item.(binding.String))
		},
	)
	//resultBox := //container.NewVBox(resultList) container.NewVScroll(resultList)
	//resultBox.Resize
	//)

	//divList := widget.NewList()
	boxList := container.NewVBox(
		searchLabel,
		searchBox,
		webLabel,
		//webBox,
		webBox2,
		resultLabel,
		//resultBox,
	)

	w.SetContent(container.NewBorder(boxList, nil, nil, nil, resultList))
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}
