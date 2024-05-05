package infra

import (
	_ "embed"
	"fyne.io/fyne/v2"
)

type AppIcon struct {
}

//go:embed "icon_tincture.png"
var icon []byte

func (a *AppIcon) AsResource() fyne.Resource {
	return fyne.NewStaticResource("icon", icon)
}
