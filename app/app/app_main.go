package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
)

type ui struct {
	window   fyne.Window
	receipts *widget.Accordion
	times    *fyne.Container
	tabs     []*container.TabItem
	tabsCont *container.AppTabs
}

var curUi *ui

func CreateWindow() fyne.Window {

	myWindow := tApp.app.NewWindow("Настойки")

	err := tApp.storage.InitDb(tApp.app.Storage().RootURI())
	if err != nil {
		log.Println(err)
	}
	return myWindow
}

func InitialLayout(w fyne.Window) fyne.Window {
	ui := getUi()
	ui.window = w

	ui.tabsCont = container.NewAppTabs(
		ui.tabs...,
	)

	w.SetContent(ui.tabsCont)

	return w
}

func getUi() *ui {
	if nil == curUi {
		receipts := makeReceipts()
		times := makeTimes()
		tabs := []*container.TabItem{
			container.NewTabItem("Мои настойки", times),
			container.NewTabItem("Рецепты", receipts),
		}
		curUi = &ui{receipts: receipts, times: times, tabs: tabs}
	}
	return curUi
}

func makeReceipts() *widget.Accordion {
	items := renderReceipts(tApp.receiptsRepository)
	accord := widget.NewAccordion(
		items...,
	)
	return accord
}

func makeTimes() *fyne.Container {
	items := tApp.tincturesRepository.GetTinctures()
	return renderTinctures(items)
}
