package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/satmaelstorm/tincture/app/domain"
	"time"
)

func renderTinctures(tinctures []domain.Tincture) *fyne.Container {
	cont := container.New(layout.NewVBoxLayout())
	for _, tincture := range tinctures {
		cont.Add(renderTincture(tincture))
	}
	return cont
}

func renderTincture(tincture domain.Tincture) *fyne.Container {
	now := time.Now()
	cont := container.New(layout.NewVBoxLayout())

	firstRow := container.New(layout.NewHBoxLayout())
	firstRow.Add(widget.NewLabelWithStyle(tincture.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
	firstRow.Add(widget.NewLabelWithStyle("Заложено\n"+tincture.CreatedAt.Format(time.DateOnly), fyne.TextAlignCenter, fyne.TextStyle{}))
	bottled := widget.NewLabelWithStyle("Переливать\n"+tincture.NeedBottledAt.Format(time.DateOnly), fyne.TextAlignCenter, fyne.TextStyle{})
	if tincture.NeedToBottled(now) {
		bottled.Importance = widget.DangerImportance
	}
	firstRow.Add(bottled)
	firstRow.Add(widget.NewLabelWithStyle("Готово к употреблению\n"+tincture.ReadyAt.Format(time.DateOnly), fyne.TextAlignCenter, fyne.TextStyle{}))
	expiredLabel := widget.NewLabelWithStyle("Срок годности\n"+tincture.ExpiredAt.Format(time.DateOnly), fyne.TextAlignCenter, fyne.TextStyle{})
	if tincture.IsExpire(now) {
		expiredLabel.Importance = widget.DangerImportance
	}
	firstRow.Add(expiredLabel)
	cont.Add(firstRow)

	progressPcts := []float64{
		tincture.ReadyToBottle(time.Now()),
		tincture.ReadyToDrink(time.Now()),
	}

	progressNames := []string{
		"Готовность перелить",
		"Готовность к употреблению",
	}

	for i, progress := range progressPcts {
		progressRow := container.New(layout.NewGridLayoutWithColumns(2))
		progressRow.Add(widget.NewLabel(progressNames[i]))
		bar := widget.NewProgressBar()
		bar.SetValue(progress)
		progressRow.Add(bar)
		cont.Add(progressRow)
	}

	return cont
}
