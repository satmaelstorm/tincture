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
	)
	w := app.CreateWindow()
	app.InitialLayout(w)
	w.ShowAndRun()
}
