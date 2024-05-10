package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
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
	myWindow.SetMaster()

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

	w.SetOnClosed(func() {
		pref := thisApp().app.Preferences()
		pref.SetBool("ready", false)
		pref.SetFloat("width", float64(w.Canvas().Size().Width))
		pref.SetFloat("height", float64(w.Canvas().Size().Height))
		//log.Printf("%v", w.Canvas().Size())
		//log.Printf("%v", thisApp().app.Settings().Scale())
		//log.Printf("%v", theme.Padding())
	})

	w.SetContent(ui.tabsCont)

	if thisApp().app.Driver().Device().IsMobile() {
		w.SetFullScreen(true)
	} else {
		resizeDesktop(ui, w)
	}

	return w
}

func resizeDesktop(ui *ui, w fyne.Window) {

	pref := thisApp().app.Preferences()
	if !pref.Bool("ready") {
		w.Resize(fyne.NewSize(ui.tabsCont.Size().Width+theme.Padding()*2, ui.tabsCont.Size().Height*1.5+theme.Padding()*2))
		return
	}
	w.Resize(fyne.NewSize(ui.tabsCont.Size().Width+theme.Padding()*2, ui.tabsCont.Size().Height*1.5+theme.Padding()*2))
	//не выходит почему-то, окно с кажды открытием растет и растет. не могу уловить закономерность. и это не Padding точно
	//именно поэтому пока ready=false
	//w.Resize(fyne.NewSize(float32(pref.Float("width")), float32(pref.Float("height"))))
}

func getUi() *ui {
	if nil == curUi {
		curUi = &ui{}
		curUi.receipts = makeReceipts()
		prepareTinctures := makePrepareTinctures()
		tabTimes := container.NewVScroll(prepareTinctures)
		tabTimes.SetMinSize(fyne.NewSize(360, 250))
		readyTinctures := makeReadyTinctures()
		curUi.tabs = []*container.TabItem{
			container.NewTabItem("Настаивается", tabTimes),
			container.NewTabItem("Погребок", container.NewVScroll(readyTinctures)),
			container.NewTabItem("Рецепты", container.NewVScroll(curUi.receipts)),
		}
	}
	return curUi
}

func makeReceipts() *widget.Accordion {
	r := &receiptRenderer{receiptsRepository: thisApp().receiptsRepository}
	return r.renderReceipts()
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
