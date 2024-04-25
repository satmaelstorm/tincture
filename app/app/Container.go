package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/satmaelstorm/tincture/app/app/port"
)

var tApp *appContainer

type appContainer struct {
	app                 fyne.App
	storage             port.InnerStorage
	receiptsRepository  port.ReceiptStorage
	tincturesRepository port.TinctureStorage
}

func InitApp(
	storage port.InnerStorage,
	receiptsStorage port.ReceiptStorage,
	tinctureStorage port.TinctureStorage,
) {
	tApp = &appContainer{
		app:                 app.NewWithID("xyz.satmaelstorm.tincture"),
		storage:             storage,
		receiptsRepository:  receiptsStorage,
		tincturesRepository: tinctureStorage,
	}
}
