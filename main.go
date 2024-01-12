/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     music-grabber
 * @Date:        2024-01-13 00:26
 * @Description:
 */

package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/leafney/music-grabber/pkg/cdp"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello")
	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hi", func() {
			//hello.SetText("Welcome!")
			cdp.StartBrowser()
		}),
	))

	w.ShowAndRun()
}
