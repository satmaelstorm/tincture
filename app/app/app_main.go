package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
)

type uiRenderers struct {
	prepareTinctureRenderer *prepareTinctureRenderer
	readyTinctureRenderer   *readyTinctureRenderer
}

type ui struct {
	window   fyne.Window
	receipts *widget.Accordion
	tabs     []*container.TabItem
	tabsCont *container.AppTabs
	render   uiRenderers
}

var curUi *ui

func CreateWindow() fyne.Window {

	myWindow := thisApp().app.NewWindow("Tinctures")

	err := thisApp().storage.InitDb(thisApp().app.Storage().RootURI())
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

	if !thisApp().app.Driver().Device().IsMobile() {
		ui.tabsCont.Resize(w.Canvas().Size())
	}

	w.SetContent(ui.tabsCont)

	return w
}

func getUi() *ui {
	if nil == curUi {
		curUi = &ui{}
		curUi.receipts = makeReceipts()
		prepareTinctures := makePrepareTinctures()
		tabTimes := container.NewVScroll(prepareTinctures)
		tabTimes.SetMinSize(fyne.NewSize(360, 250))
		readyTinctures := container.NewVScroll(makeReadyTinctures())
		curUi.tabs = []*container.TabItem{
			container.NewTabItem("Настаивается", tabTimes),
			container.NewTabItem("Погребок", readyTinctures),
			container.NewTabItem("Рецепты", curUi.receipts),
		}
	}
	return curUi
}

func makeReceipts() *widget.Accordion {
	items := renderReceipts(thisApp().receiptsRepository)
	accord := widget.NewAccordion(
		items...,
	)
	return accord
}

func makePrepareTinctures() *fyne.Container {
	renderer := &prepareTinctureRenderer{tinctureRepository: thisApp().tincturesRepository}
	curUi.render.prepareTinctureRenderer = renderer
	items := thisApp().tincturesRepository.GetPreparingTinctures()
	return renderer.renderTinctures(items)
}

func makeReadyTinctures() *fyne.Container {
	renderer := &readyTinctureRenderer{tinctureRepository: thisApp().tincturesRepository}
	curUi.render.readyTinctureRenderer = renderer
	items := thisApp().tincturesRepository.GetReadyTinctures()
	return renderer.renderTinctures(items)
}
