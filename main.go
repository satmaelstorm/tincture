package main

import (
	"log"

	"github.com/satmaelstorm/tincture/app/app"
	"github.com/satmaelstorm/tincture/app/infra"
)

func main() {
	// Initialize localization
	i18n, err := infra.NewI18nAdapter()
	if err != nil {
		log.Fatal(err)
	}

	db := new(infra.TinctureDB)
	app.InitApp(
		db,
		db,
		db,
		new(infra.AppIcon),
		new(infra.Icons),
		infra.NewDispatcher(),
		i18n,
	)
	w := app.CreateWindow()
	app.InitialLayout(w)
	w.ShowAndRun()
}
