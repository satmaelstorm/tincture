package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/satmaelstorm/tincture/app/app/port"
	"github.com/satmaelstorm/tincture/app/domain"
	"sort"
	"time"
)

type readyTinctureRenderer struct {
	tinctureRepository port.TinctureStorage
	cont               *fyne.Container
	rows               map[string]*fyne.Container
	tinctures          map[string]domain.Tincture
}

func (r *readyTinctureRenderer) addTincture(tincture domain.Tincture) {
	r.addRenderTincture(tincture)
	r.rearrangeRows()
	r.cont.Refresh()
}

func (r *readyTinctureRenderer) renderTinctures(tinctures []domain.Tincture) *fyne.Container {
	r.rows = make(map[string]*fyne.Container, len(tinctures))
	r.tinctures = make(map[string]domain.Tincture, len(tinctures))
	r.cont = container.New(layout.NewVBoxLayout())
	for _, tincture := range tinctures {
		r.addRenderTincture(tincture)
	}
	return r.cont
}

func (r *readyTinctureRenderer) addRenderTincture(tincture domain.Tincture) {
	rowCont := container.New(layout.NewVBoxLayout())
	r.rows[tincture.Uuid.String()] = r.renderTincture(tincture, rowCont)
	r.cont.Add(rowCont)
	r.tinctures[tincture.Uuid.String()] = tincture
}

func (r *readyTinctureRenderer) renderTincture(tincture domain.Tincture, cont *fyne.Container) *fyne.Container {
	now := time.Now()
	cont.RemoveAll()
	firstRow := container.New(layout.NewAdaptiveGridLayout(4))
	title := widget.NewLabelWithStyle(tincture.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	title.Wrapping = fyne.TextWrapWord
	firstRow.Add(title)
	readiness := widget.NewLabelWithStyle("Готово к употреблению\n"+tincture.ReadyAt.Format(time.DateOnly), fyne.TextAlignCenter, fyne.TextStyle{})
	if tincture.IsReady(now) {
		readiness.Importance = widget.SuccessImportance
	}
	firstRow.Add(readiness)
	expiredLabel := widget.NewLabelWithStyle("Срок годности\n"+tincture.ExpiredAt.Format(time.DateOnly), fyne.TextAlignCenter, fyne.TextStyle{})
	if tincture.IsNearExpire(now) {
		expiredLabel.Importance = widget.DangerImportance
	}
	firstRow.Add(expiredLabel)

	firstRow = r.addTinctureControlPanel(tincture, firstRow)

	cont.Add(firstRow)

	return cont
}

func (r *readyTinctureRenderer) addTinctureControlPanel(tincture domain.Tincture, cont *fyne.Container) *fyne.Container {

	drank := widget.NewButton("Выпито", func() {
		r.handleDeletedButton(tincture)
	})
	drank.Importance = widget.DangerImportance

	cont.Add(drank)
	return cont
}

func (r *readyTinctureRenderer) handleDeletedButton(tincture domain.Tincture) {
	r.cont.Remove(r.rows[tincture.Uuid.String()])
	delete(r.rows, tincture.Uuid.String())
	r.tinctureRepository.DeleteTincture(&tincture)
	delete(r.tinctures, tincture.Uuid.String())
}

func (r *readyTinctureRenderer) rearrangeRows() {
	tinctures := make([]domain.Tincture, 0, len(r.tinctures))
	for _, tincture := range r.tinctures {
		tinctures = append(tinctures, tincture)
	}
	sort.Slice(tinctures, func(i, j int) bool {
		if tinctures[i].ReadyAt.Before(tinctures[j].ReadyAt) {
			return true
		}
		if tinctures[i].ExpiredAt.Before(tinctures[j].ExpiredAt) {
			return true
		}
		return false
	})
	for _, obj := range r.rows {
		r.cont.Remove(obj)

	}
	for _, tincture := range tinctures {
		r.cont.Add(r.rows[tincture.Uuid.String()])
	}
}