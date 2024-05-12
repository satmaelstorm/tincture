package main

import (
	"github.com/satmaelstorm/tincture/app/app"
	"github.com/satmaelstorm/tincture/app/infra"
)

func main() {
	db := new(infra.TinctureDB)
	app.InitApp(
		db,
		db,
		db,
		new(infra.AppIcon),
		new(infra.Icons),
		infra.NewDispatcher(),
	)
	w := app.CreateWindow()
	app.InitialLayout(w)
	w.ShowAndRun()
}
