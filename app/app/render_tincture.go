package app

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	widgetX "fyne.io/x/fyne/widget"
	"github.com/satmaelstorm/tincture/app/app/port"
	"github.com/satmaelstorm/tincture/app/domain"
	"sort"
	"strconv"
	"time"
)

type tinctureRenderer struct {
	tinctureRepository port.TinctureStorage
	cont               *fyne.Container
	rows               map[string]*fyne.Container
	addPopup           *widget.PopUp
	tinctures          map[string]domain.Tincture
}

func (t *tinctureRenderer) renderTinctures(tinctures []domain.Tincture) *fyne.Container {
	t.rows = make(map[string]*fyne.Container, len(tinctures))
	t.tinctures = make(map[string]domain.Tincture, len(tinctures))
	t.cont = container.New(layout.NewVBoxLayout())
	t.cont.Add(widget.NewButton("Добавить", func() {
		t.handleAddButton()
	}))
	for _, tincture := range tinctures {
		t.addRenderTincture(tincture)
	}
	return t.cont
}

func (t *tinctureRenderer) addRenderTincture(tincture domain.Tincture) {
	rowCont := container.New(layout.NewVBoxLayout())
	t.rows[tincture.Uuid.String()] = t.renderTincture(tincture, rowCont)
	t.cont.Add(rowCont)
	t.tinctures[tincture.Uuid.String()] = tincture
}

func (t *tinctureRenderer) renderTincture(tincture domain.Tincture, cont *fyne.Container) *fyne.Container {
	now := time.Now()
	cont.RemoveAll()
	firstRow := container.New(layout.NewAdaptiveGridLayout(7))
	title := widget.NewLabelWithStyle(tincture.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	title.Wrapping = fyne.TextWrapWord
	firstRow.Add(title)
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

	firstRow = t.addTinctureControlPanel(tincture, firstRow)

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
		progressRow := container.New(layout.NewAdaptiveGridLayout(2))
		progressRow.Add(widget.NewLabel(progressNames[i]))
		bar := widget.NewProgressBar()
		bar.SetValue(progress)
		progressRow.Add(bar)
		cont.Add(progressRow)
	}

	return cont
}

func (t *tinctureRenderer) addTinctureControlPanel(tincture domain.Tincture, cont *fyne.Container) *fyne.Container {
	bottled := widget.NewButton("Перелито", func() {
		t.handleBottledButton(tincture)
	})
	if tincture.NeedToBottled(time.Now()) {
		bottled.Importance = widget.DangerImportance
	} else if tincture.IsBottled() {
		bottled.Importance = widget.SuccessImportance
		bottled.Disable()
	} else {
		bottled.Importance = widget.HighImportance
	}

	drank := widget.NewButton("Выпито", func() {
		t.handleDeletedButton(tincture)
	})
	drank.Importance = widget.DangerImportance

	cont.Add(bottled)
	cont.Add(drank)
	return cont
}

func (t *tinctureRenderer) handleBottledButton(tincture domain.Tincture) {
	tincture.Bottled(time.Now())
	t.tinctureRepository.SaveTincture(&tincture)
	cont := t.rows[tincture.Uuid.String()]
	t.renderTincture(tincture, cont)
	cont.Refresh()
}

func (t *tinctureRenderer) handleDeletedButton(tincture domain.Tincture) {
	t.cont.Remove(t.rows[tincture.Uuid.String()])
	delete(t.rows, tincture.Uuid.String())
	t.tinctureRepository.DeleteTincture(&tincture)
}

func (t *tinctureRenderer) handleAddButton() {
	t.createAddPopup()
	t.addPopup.Show()
}

func (t *tinctureRenderer) createAddPopup() {
	if nil != t.addPopup {
		return
	}
	titleEntry := widget.NewEntry()
	titleEntry.Validator = func(s string) error {
		if len(s) < 3 || len(s) > 120 {
			return errors.New("название должно быть от 3 до 120 символов")
		}
		return nil
	}
	titleEntry.PlaceHolder = "Название от 3 до 120 символов"
	createdAtEntry := widget.NewEntry()
	createdAtEntry.SetText(time.Now().Format(time.DateOnly))
	createdAtEntry.Validator = func(s string) error {
		_, err := time.Parse(time.DateOnly, s)
		return err
	}
	daysToBottle := widget.NewEntry()
	daysToBottle.SetText("14")
	daysToBottle.Validator = func(s string) error {
		i, err := strconv.Atoi(s)
		if i < 0 {
			return errors.New("число не может быть меньше 0")
		}
		return err
	}
	daysToRest := widget.NewEntry()
	daysToRest.SetText("0")
	daysToRest.Validator = func(s string) error {
		i, err := strconv.Atoi(s)
		if i < 0 {
			return errors.New("число не может быть меньше 0")
		}
		return err
	}
	daysToExpire := widget.NewEntry()
	daysToExpire.SetText("365")
	daysToExpire.Validator = func(s string) error {
		i, err := strconv.Atoi(s)
		if i < 0 {
			return errors.New("число не может быть меньше 0")
		}
		return err
	}

	calendarPopup := widget.NewModalPopUp(widget.NewEntry(), getUi().window.Canvas())

	createdAtCalendar := widgetX.NewCalendar(time.Now(), func(t time.Time) {
		createdAtEntry.SetText(t.Format(time.DateOnly))
		calendarPopup.Hide()
	})
	calendarPopup.Content = createdAtCalendar

	createdAtRow := container.New(layout.NewGridLayoutWithColumns(2), createdAtEntry, widget.NewButton("Выбрать день", func() {
		calendarPopup.Show()
	}))

	form := widget.NewForm(
		widget.NewFormItem("Название:", titleEntry),
		widget.NewFormItem("Заложено:", createdAtRow),
		widget.NewFormItem("Настаивать дней:", daysToBottle),
		widget.NewFormItem("Отдыхать дней:", daysToRest),
		widget.NewFormItem("Годно дней:", daysToExpire),
	)

	popup := widget.NewModalPopUp(widget.NewCard("Добавить настойку", "", form), getUi().window.Canvas())

	form.OnSubmit = func() {
		tincture := domain.NewTincture(
			titleEntry.Text,
			createdAtEntry.Text,
			daysToBottle.Text,
			daysToRest.Text,
			daysToExpire.Text,
		)
		t.handleSubmitNewTincture(tincture)
		popup.Hide()
	}

	form.OnCancel = func() {
		popup.Hide()
	}

	popup.Resize(fyne.NewSize(t.cont.Size().Width*0.9, t.cont.Size().Height*0.9))
	t.addPopup = popup
}

func (t *tinctureRenderer) handleSubmitNewTincture(tincture domain.Tincture) {
	t.tinctureRepository.CreateTincture(&tincture)
	t.addRenderTincture(tincture)
	t.rearrangeRows()
	t.cont.Refresh()
}

func (t *tinctureRenderer) rearrangeRows() {
	tinctures := make([]domain.Tincture, 0, len(t.tinctures))
	for _, tincture := range t.tinctures {
		tinctures = append(tinctures, tincture)
	}
	sort.Slice(tinctures, func(i, j int) bool {
		if tinctures[i].NeedBottledAt.Before(tinctures[j].NeedBottledAt) {
			return true
		}
		if tinctures[i].CreatedAt.Before(tinctures[j].CreatedAt) {
			return true
		}
		return false
	})
	for _, obj := range t.rows {
		t.cont.Remove(obj)

	}
	for _, tincture := range tinctures {
		t.cont.Add(t.rows[tincture.Uuid.String()])
	}
}
