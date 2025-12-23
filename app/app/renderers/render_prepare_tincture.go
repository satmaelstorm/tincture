package renderers

import (
	"errors"
	"sort"
	"strconv"
	"time"

	"github.com/satmaelstorm/tincture/app/domain/events"

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
	bus                port.EventBus
	localizer          port.Localizer

	cont      *fyne.Container
	rows      map[string]*fyne.Container
	addPopup  *widget.PopUp
	tinctures map[string]domain.Tincture

	formTitleEntry *widget.Entry
}

func NewPrepareTinctureRenderer(
	repository port.TinctureStorage,
	canvas fyne.Canvas,
	bus port.EventBus,
	localizer port.Localizer,
) *PrepareTinctureRenderer {
	return &PrepareTinctureRenderer{
		tinctureRepository: repository,
		canvas:             canvas,
		bus:                bus,
		localizer:          localizer,
	}
}

func (p *PrepareTinctureRenderer) RenderTinctures(tinctures []domain.Tincture) *fyne.Container {
	p.rows = make(map[string]*fyne.Container, len(tinctures))
	p.tinctures = make(map[string]domain.Tincture, len(tinctures))
	p.cont = container.New(layout.NewVBoxLayout())
	p.cont.Add(container.NewPadded(widget.NewButton(p.localizer.Localize("addTinctureButton"), func() {
		p.bus.Dispatch(&events.TinctureAddButton{})
	})))
	for _, tincture := range tinctures {
		p.addRenderTincture(tincture)
	}
	return p.cont
}

func (p *PrepareTinctureRenderer) ShowAddPopup() {
	p.createAddPopup()
	p.formTitleEntry.SetText("")
	p.addPopup.Show()
}

func (p *PrepareTinctureRenderer) HideAddPopup() {
	p.addPopup.Hide()
}

func (p *PrepareTinctureRenderer) AddTincture(tincture domain.Tincture) {
	p.addRenderTincture(tincture)
	p.rearrangeRows()
	p.cont.Refresh()
}

func (p *PrepareTinctureRenderer) RemoveTincture(tincture domain.Tincture) {
	oldRow := p.rows[tincture.Uuid.String()]
	delete(p.rows, tincture.Uuid.String())
	delete(p.tinctures, tincture.Uuid.String())
	p.cont.Remove(oldRow)
	p.cont.Refresh()
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
	firstRow.Add(widget.NewLabelWithStyle(p.localizer.Localize("tinctureCreatedAt")+"\n"+tincture.CreatedAt.Format(time.DateOnly), fyne.TextAlignCenter, fyne.TextStyle{}))
	bottled := widget.NewLabelWithStyle(p.localizer.Localize("tinctureNeedBottledAt")+"\n"+tincture.NeedBottledAt.Format(time.DateOnly), fyne.TextAlignCenter, fyne.TextStyle{})
	if tincture.NeedToBottled(now) {
		bottled.Importance = widget.DangerImportance
	}
	firstRow.Add(bottled)

	firstRow = p.addTinctureControlPanel(tincture, firstRow)

	cont.Add(firstRow)

	progressRow := container.New(layout.NewAdaptiveGridLayout(2))
	progressRow.Add(widget.NewLabel(p.localizer.Localize("tinctureBottleProgress")))
	bar := widget.NewProgressBar()
	bar.SetValue(tincture.ReadyToBottle(now))
	progressRow.Add(bar)
	cont.Add(progressRow)

	return cont
}

func (p *PrepareTinctureRenderer) addTinctureControlPanel(tincture domain.Tincture, cont *fyne.Container) *fyne.Container {
	bottled := widget.NewButton(p.localizer.Localize("bottleTinctureButton"), func() {
		p.bus.Dispatch(&events.TinctureBottled{Tincture: tincture})
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
			return errors.New(p.localizer.Localize("tinctureNameValidationError"))
		}
		return nil
	}
	titleEntry.PlaceHolder = p.localizer.Localize("tinctureNamePlaceholder")
	p.formTitleEntry = titleEntry

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
			return errors.New(p.localizer.Localize("tinctureDaysValidationError"))
		}
		return err
	}
	daysToRest := widget.NewEntry()
	daysToRest.SetText("0")
	daysToRest.Validator = func(s string) error {
		i, err := strconv.Atoi(s)
		if i < 0 {
			return errors.New(p.localizer.Localize("tinctureDaysValidationError"))
		}
		return err
	}
	daysToExpire := widget.NewEntry()
	daysToExpire.SetText("365")
	daysToExpire.Validator = func(s string) error {
		i, err := strconv.Atoi(s)
		if i < 0 {
			return errors.New(p.localizer.Localize("tinctureDaysValidationError"))
		}
		return err
	}

	calendarPopup := widget.NewModalPopUp(widget.NewEntry(), p.canvas)

	createdAtCalendar := widgetX.NewCalendar(time.Now(), func(t time.Time) {
		createdAtEntry.SetText(t.Format(time.DateOnly))
		calendarPopup.Hide()
	})
	calendarPopup.Content = createdAtCalendar

	createdAtRow := container.New(layout.NewAdaptiveGridLayout(2), createdAtEntry, widget.NewButton(p.localizer.Localize("tinctureChooseDay"), func() {
		calendarPopup.Show()
	}))

	form := widget.NewForm(
		widget.NewFormItem(p.localizer.Localize("tinctureName")+":", titleEntry),
		widget.NewFormItem(p.localizer.Localize("tinctureCreatedAt")+":", createdAtRow),
		widget.NewFormItem(p.localizer.Localize("tinctureDaysToBottle")+":", daysToBottle),
		widget.NewFormItem(p.localizer.Localize("tinctureDaysToRest")+":", daysToRest),
		widget.NewFormItem(p.localizer.Localize("tinctureDaysToExpire")+":", daysToExpire),
	)

	popup := widget.NewModalPopUp(widget.NewCard(p.localizer.Localize("tinctureAddPopupTitle"), "", form), p.canvas)

	form.OnSubmit = func() {
		tincture := domain.NewTincture(
			titleEntry.Text,
			createdAtEntry.Text,
			daysToBottle.Text,
			daysToRest.Text,
			daysToExpire.Text,
		)
		p.bus.Dispatch(&events.TinctureSubmit{Tincture: tincture})
	}

	form.OnCancel = func() {
		p.bus.Dispatch(new(events.TinctureCancelButton))
	}

	form.CancelText = p.localizer.Localize("tinctureAddPopupCancel")
	form.SubmitText = p.localizer.Localize("tinctureAddPopupSubmit")

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
