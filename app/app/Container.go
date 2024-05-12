package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/satmaelstorm/tincture/app/app/handlers"
	"github.com/satmaelstorm/tincture/app/app/port"
	"github.com/satmaelstorm/tincture/app/app/renderers"
)

var tApp *appContainer

type appContainer struct {
	app                 fyne.App
	storage             port.InnerStorage
	receiptsRepository  port.ReceiptStorage
	tincturesRepository port.TinctureStorage
	icons               port.AppIcons
	dispatcher          port.EventDispatcher

	handlerReceiptForm *handlers.ReceiptFormHandlers
}

func InitApp(
	storage port.InnerStorage,
	receiptsStorage port.ReceiptStorage,
	tinctureStorage port.TinctureStorage,
	appIconProvider port.AppIcon,
	iconsProvider port.AppIcons,
	dispatcher port.EventDispatcher,
) {
	tApp = &appContainer{
		app:                 app.NewWithID("xyz.satmaelstorm.tincture"),
		storage:             storage,
		receiptsRepository:  receiptsStorage,
		tincturesRepository: tinctureStorage,
		icons:               iconsProvider,
		dispatcher:          dispatcher,
	}
	tApp.app.SetIcon(appIconProvider.AsResource())
}

func (a *appContainer) initReceiptHandlers(
	r *renderers.ReceiptRenderer,
	f *renderers.ReceiptEditForm,
) {
	handler := handlers.NewReceiptFormHandlers(r, f, a.receiptsRepository)
	a.dispatcher.AddSubscriber(handler)
	a.handlerReceiptForm = handler
}

func thisApp() *appContainer {
	return tApp
}
