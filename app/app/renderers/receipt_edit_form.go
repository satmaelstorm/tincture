package renderers

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/satmaelstorm/tincture/app/app/port"
	"github.com/satmaelstorm/tincture/app/domain"
	"github.com/satmaelstorm/tincture/app/domain/events"
)

type ReceiptEditForm struct {
	form   *widget.Form
	title  *widget.Entry
	desc   *widget.Entry
	items  *fyne.Container
	popup  *widget.PopUp
	canvas fyne.Canvas

	bus port.EventBus

	currentReceipt *domain.Receipt
}

func NewReceiptEditForm(
	bus port.EventBus,
	canvas fyne.Canvas,
) *ReceiptEditForm {
	r := new(ReceiptEditForm)
	r.bus = bus
	r.canvas = canvas

	titleEntry := widget.NewEntry()
	titleEntry.PlaceHolder = "От 3 букв"
	titleEntry.Validator = func(s string) error {
		l := len([]rune(s))
		if l < 3 || l > 120 {
			return errors.New("название должно быть от 3 до 120 символов")
		}
		return nil
	}

	descEntry := widget.NewMultiLineEntry()
	descEntry.SetMinRowsVisible(5)
	descEntry.Wrapping = fyne.TextWrapWord
	descEntry.Validator = func(s string) error {
		l := len([]rune(s))
		if l > 2000 {
			return errors.New("описание должно быть от 0 до 2000 символов")
		}
		return nil
	}
	descEntry.PlaceHolder = "От 0 до 2000 букв"

	itemsWithButton := container.NewVBox()
	itemsCont := container.NewAdaptiveGrid(2)
	itemsWithButton.Add(widget.NewButton("Добавить", func() {
		r.addNewItemRow(false, domain.ReceiptItem{})
	}))
	itemsWithButton.Add(itemsCont)

	form := widget.NewForm(
		widget.NewFormItem("Название", titleEntry),
		widget.NewFormItem("Описание", descEntry),
		widget.NewFormItem("Ингридиенты", itemsWithButton),
	)

	form.CancelText = "Отмена"
	form.SubmitText = "Ок"

	form.OnCancel = func() {
		r.bus.Dispatch(new(events.ReceiptFormCancel))
	}

	form.OnSubmit = func() {
		r.bus.Dispatch(new(events.ReceiptFormSubmit))
	}

	card := widget.NewCard("Рецепт", "", form)

	popup := widget.NewModalPopUp(container.NewVScroll(card), r.canvas)

	r.form = form
	r.title = titleEntry
	r.desc = descEntry
	r.items = itemsCont
	r.popup = popup

	return r
}

func (r *ReceiptEditForm) Clear() {
	r.desc.Text = ""
	r.title.Text = ""
	r.currentReceipt = nil
	r.items.RemoveAll()
}

func (r *ReceiptEditForm) FromReceipt(receipt domain.Receipt) {
	r.desc.Text = receipt.Description
	r.title.Text = receipt.Title
	r.currentReceipt = &receipt
	r.items.RemoveAll()
	for _, item := range receipt.Items {
		r.addNewItemRow(true, item)
	}
}

func (r *ReceiptEditForm) addNewItemRow(addItem bool, item domain.ReceiptItem) {
	nEntry := widget.NewEntry()
	nEntry.PlaceHolder = "Название"
	cEntry := widget.NewEntry()
	cEntry.PlaceHolder = "Кол-во"
	if addItem {
		nEntry.Text = item.Name
		cEntry.Text = item.Quantity
	}
	r.items.Add(nEntry)
	r.items.Add(cEntry)
}

func (r *ReceiptEditForm) Hide() {
	r.popup.Hide()
}

func (r *ReceiptEditForm) Show() {
	r.popup.Resize(r.canvas.Size())
	r.popup.Show()
}

func (r *ReceiptEditForm) CollectReceipt() (domain.Receipt, bool) {
	l := len(r.items.Objects)
	items := make([]domain.ReceiptItem, 0, l/2)
	for i := 0; i < l; i += 2 {
		items = append(items, domain.NewReceiptItem(r.items.Objects[i].(*widget.Entry).Text, r.items.Objects[i+1].(*widget.Entry).Text))
	}
	if nil == r.currentReceipt {
		return domain.NewReceipt(r.title.Text, r.desc.Text, items...), true
	}
	r.currentReceipt.Modify(r.title.Text, r.desc.Text, items...)
	return *r.currentReceipt, false
}
