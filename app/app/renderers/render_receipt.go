package renderers

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/satmaelstorm/tincture/app/app/port"
	"github.com/satmaelstorm/tincture/app/domain"
	"github.com/satmaelstorm/tincture/app/domain/events"
)

type ReceiptRenderer struct {
	receiptsRepository port.ReceiptStorage
	accord             *widget.Accordion
	cont               *fyne.Container
	accItems           *receiptAccordionItems
	bus                port.EventBus
}

func NewReceiptRenderer(
	repository port.ReceiptStorage,
	bus port.EventBus,
) *ReceiptRenderer {
	return &ReceiptRenderer{receiptsRepository: repository, bus: bus}
}

func (r *ReceiptRenderer) RenderReceipts() *fyne.Container {
	receipts := r.receiptsRepository.GetReceipts()
	aItems := make([]*widget.AccordionItem, 0, len(receipts))
	r.accItems = newReceiptAccordionItems(len(receipts))
	for _, receipt := range receipts {
		aItems = append(aItems, r.renderReceipt(receipt, r.accItems.makeNew(receipt)))
	}
	r.cont = container.NewVBox()
	r.cont.Add(widget.NewButton("Добавить", func() {
		r.bus.Dispatch(new(events.ReceiptAddButton))
	}))
	r.accord = widget.NewAccordion(aItems...)
	r.cont.Add(r.accord)
	return r.cont
}

func (r *ReceiptRenderer) AddReceipt(receipt domain.Receipt) {
	item := r.accItems.makeNew(receipt)
	item = r.renderReceipt(receipt, item)
	r.accord.Append(item)
	r.cont.Refresh()
}

func (r *ReceiptRenderer) RefreshReceipt(receipt domain.Receipt) {
	item := r.accItems.get(receipt)
	item = r.renderReceipt(receipt, item)
	r.accord.Refresh()
}

func (r *ReceiptRenderer) RemoveReceipt(receipt domain.Receipt) {
	item := r.accItems.get(receipt)
	r.accord.Remove(item)
	r.accord.Refresh()
}

func (r *ReceiptRenderer) renderReceipt(receipt domain.Receipt, accItem *widget.AccordionItem) *widget.AccordionItem {
	receiptContainer := container.New(layout.NewAdaptiveGridLayout(2))

	textDesc := widget.NewRichText(&widget.TextSegment{
		Style: widget.RichTextStyleParagraph,
		Text:  receipt.Description,
	})
	textDesc.Wrapping = fyne.TextWrapWord
	desc := container.NewScroll(textDesc)
	receiptContainer.Add(desc)
	itemsContainer := container.New(layout.NewVBoxLayout())
	for _, item := range receipt.Items {
		lName := widget.NewLabel(item.Name)
		lName.Wrapping = fyne.TextWrapWord
		itemsContainer.Add(container.New(
			layout.NewGridLayoutWithColumns(2),
			lName,
			widget.NewLabel(item.Quantity),
		),
		)
	}
	receiptContainer.Add(itemsContainer)

	editButton := widget.NewButton("Редактировать", func() {
		r.bus.Dispatch(&events.ReceiptEditButton{ReceiptUuid: receipt.Uuid})
	})

	delButton := widget.NewButton("Удалить", func() {
		r.bus.Dispatch(&events.ReceiptDeleteButton{ReceiptUuid: receipt.Uuid})
	})
	delButton.Importance = widget.DangerImportance

	buttons := container.NewBorder(widget.NewLabel(""), widget.NewLabel(""), editButton, delButton)

	accItem.Title = receipt.Title
	accItem.Detail = container.NewVBox(receiptContainer, buttons)
	return accItem
}
