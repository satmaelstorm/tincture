package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"github.com/satmaelstorm/tincture/app/app/renderers"
	"log"
)

type uiRenderers struct {
	prepareTinctureRenderer *renderers.PrepareTinctureRenderer
	readyTinctureRenderer   *renderers.ReadyTinctureRenderer
}

type ui struct {
	window   fyne.Window
	receipts *fyne.Container
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
	ui := makeUi(w)

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
		//w.SetFullScreen(true)
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

func makeUi(w fyne.Window) *ui {
	if nil == curUi {
		curUi = &ui{window: w}
		curUi.receipts = makeReceipts(w.Canvas())
		readyTinctures := makeReadyTinctures()
		prepareTinctures := makePrepareTinctures()

		thisApp().initTinctureHandlers(curUi.render.readyTinctureRenderer, curUi.render.prepareTinctureRenderer)

		tabTimes := container.NewVScroll(prepareTinctures)
		tabTimes.SetMinSize(fyne.NewSize(360, 250))
		curUi.tabs = []*container.TabItem{
			container.NewTabItem("Настаивается", tabTimes),
			container.NewTabItem("Погребок", container.NewVScroll(readyTinctures)),
			container.NewTabItem("Рецепты", container.NewVScroll(curUi.receipts)),
		}
	}
	return curUi
}

func makeReceipts(canvas fyne.Canvas) *fyne.Container {
	r := renderers.NewReceiptRenderer(thisApp().receiptsRepository, thisApp().dispatcher)
	rCont := r.RenderReceipts()
	thisApp().initReceiptHandlers(r, renderers.NewReceiptEditForm(thisApp().dispatcher, canvas))
	return rCont
}

func makeReadyTinctures() *fyne.Container {
	renderer := renderers.NewReadyTinctureRenderer(
		thisApp().dispatcher,
	)
	curUi.render.readyTinctureRenderer = renderer
	items := thisApp().tincturesRepository.GetReadyTinctures()
	return renderer.RenderTinctures(items)
}

func makePrepareTinctures() *fyne.Container {
	renderer := renderers.NewPrepareTinctureRenderer(
		thisApp().tincturesRepository,
		curUi.window.Canvas(),
		thisApp().dispatcher,
	)
	curUi.render.prepareTinctureRenderer = renderer
	items := thisApp().tincturesRepository.GetPreparingTinctures()
	return renderer.RenderTinctures(items)
}
