package renderers

import (
	"errors"
	"sort"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	widgetX "fyne.io/x/fyne/widget"
	"github.com/satmaelstorm/tincture/app/app/port"
	"github.com/satmaelstorm/tincture/app/domain"
)

type PrepareTinctureRenderer struct {
	tinctureRepository port.TinctureStorage
	canvas             fyne.Canvas
	readyTinctures     *ReadyTinctureRenderer

	cont      *fyne.Container
	rows      map[string]*fyne.Container
	addPopup  *widget.PopUp
	tinctures map[string]domain.Tincture
}

func NewPrepareTinctureRenderer(
	repository port.TinctureStorage,
	canvas fyne.Canvas,
	renderer *ReadyTinctureRenderer,
) *PrepareTinctureRenderer {
	return &PrepareTinctureRenderer{
		tinctureRepository: repository,
		canvas:             canvas,
		readyTinctures:     renderer,
	}
}

func (p *PrepareTinctureRenderer) RenderTinctures(tinctures []domain.Tincture) *fyne.Container {
	p.rows = make(map[string]*fyne.Container, len(tinctures))
	p.tinctures = make(map[string]domain.Tincture, len(tinctures))
	p.cont = container.New(layout.NewVBoxLayout())
	p.cont.Add(widget.NewButton("Добавить", func() {
		p.handleAddButton()
	}))
	for _, tincture := range tinctures {
		p.addRenderTincture(tincture)
	}
	return p.cont
}

func (p *PrepareTinctureRenderer) addRenderTincture(tincture domain.Tincture) {
	rowCont := container.New(layout.NewVBoxLayout())
	p.rows[tincture.Uuid.String()] = p.renderTincture(tincture, rowCont)
	p.cont.Add(rowCont)
	p.tinctures[tincture.Uuid.String()] = tincture
}

func (p *PrepareTinctureRenderer) renderTincture(tincture domain.Tincture, cont *fyne.Container) *fyne.Container {
	now := time.Now()
	cont.RemoveAll()
	firstRow := container.New(layout.NewAdaptiveGridLayout(4))
	title := widget.NewLabelWithStyle(tincture.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	title.Wrapping = fyne.TextWrapWord
	firstRow.Add(title)
	firstRow.Add(widget.NewLabelWithStyle("Заложено\n"+tincture.CreatedAt.Format(time.DateOnly), fyne.TextAlignCenter, fyne.TextStyle{}))
	bottled := widget.NewLabelWithStyle("Переливать\n"+tincture.NeedBottledAt.Format(time.DateOnly), fyne.TextAlignCenter, fyne.TextStyle{})
	if tincture.NeedToBottled(now) {
		bottled.Importance = widget.DangerImportance
	}
	firstRow.Add(bottled)

	firstRow = p.addTinctureControlPanel(tincture, firstRow)

	cont.Add(firstRow)

	progressRow := container.New(layout.NewAdaptiveGridLayout(2))
	progressRow.Add(widget.NewLabel("Готовность перелить"))
	bar := widget.NewProgressBar()
	bar.SetValue(tincture.ReadyToBottle(now))
	progressRow.Add(bar)
	cont.Add(progressRow)

	return cont
}

func (p *PrepareTinctureRenderer) addTinctureControlPanel(tincture domain.Tincture, cont *fyne.Container) *fyne.Container {
	bottled := widget.NewButton("Перелито", func() {
		p.handleBottledButton(tincture)
	})
	if tincture.NeedToBottled(time.Now()) {
		bottled.Importance = widget.DangerImportance
	} else if tincture.IsBottled() {
		bottled.Importance = widget.SuccessImportance
		bottled.Disable()
	} else {
		bottled.Importance = widget.HighImportance
	}

	cont.Add(bottled)
	return cont
}

func (p *PrepareTinctureRenderer) createAddPopup() {
	if nil != p.addPopup {
		return
	}
	titleEntry := widget.NewEntry()
	titleEntry.Validator = func(s string) error {
		l := len([]rune(s))
		if l < 3 || l > 120 {
			return errors.New("название должно быть от 3 до 120 символов")
		}
		return nil
	}
	titleEntry.PlaceHolder = "От 3 букв"
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

	calendarPopup := widget.NewModalPopUp(widget.NewEntry(), p.canvas)

	createdAtCalendar := widgetX.NewCalendar(time.Now(), func(t time.Time) {
		createdAtEntry.SetText(t.Format(time.DateOnly))
		calendarPopup.Hide()
	})
	calendarPopup.Content = createdAtCalendar

	createdAtRow := container.New(layout.NewAdaptiveGridLayout(2), createdAtEntry, widget.NewButton("Выбрать день", func() {
		calendarPopup.Show()
	}))

	form := widget.NewForm(
		widget.NewFormItem("Название:", titleEntry),
		widget.NewFormItem("Заложено:", createdAtRow),
		widget.NewFormItem("Настаивать дней:", daysToBottle),
		widget.NewFormItem("Отдыхать дней:", daysToRest),
		widget.NewFormItem("Годно дней:", daysToExpire),
	)

	popup := widget.NewModalPopUp(widget.NewCard("Добавить настойку", "", form), p.canvas)

	form.OnSubmit = func() {
		tincture := domain.NewTincture(
			titleEntry.Text,
			createdAtEntry.Text,
			daysToBottle.Text,
			daysToRest.Text,
			daysToExpire.Text,
		)
		p.handleSubmitNewTincture(tincture)
		titleEntry.SetText("")
		popup.Hide()
	}

	form.OnCancel = func() {
		titleEntry.SetText("")
		popup.Hide()
	}

	form.CancelText = "Отмена"
	form.SubmitText = "Ок"

	popup.Resize(fyne.NewSize(p.cont.Size().Width*0.5, p.cont.Size().Height*0.3))
	p.addPopup = popup
}

func (p *PrepareTinctureRenderer) rearrangeRows() {
	tinctures := make([]domain.Tincture, 0, len(p.tinctures))
	for _, tincture := range p.tinctures {
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
	for _, obj := range p.rows {
		p.cont.Remove(obj)

	}
	for _, tincture := range tinctures {
		p.cont.Add(p.rows[tincture.Uuid.String()])
	}
}

func (p *PrepareTinctureRenderer) handleBottledButton(tincture domain.Tincture) {
	tincture.Bottled(time.Now())
	p.tinctureRepository.SaveTincture(&tincture)
	p.readyTinctures.addTincture(tincture)
	oldRow := p.rows[tincture.Uuid.String()]
	delete(p.rows, tincture.Uuid.String())
	delete(p.tinctures, tincture.Uuid.String())
	p.cont.Remove(oldRow)
	p.cont.Refresh()
}

func (p *PrepareTinctureRenderer) handleAddButton() {
	p.createAddPopup()
	p.addPopup.Show()
}

func (p *PrepareTinctureRenderer) handleSubmitNewTincture(tincture domain.Tincture) {
	p.tinctureRepository.CreateTincture(&tincture)
	p.addRenderTincture(tincture)
	p.rearrangeRows()
	p.cont.Refresh()
}
