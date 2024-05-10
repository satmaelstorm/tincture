package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/satmaelstorm/tincture/app/app/port"
	"github.com/satmaelstorm/tincture/app/domain"
)

type receiptEditForm struct {
	title *widget.Entry
	desc  *widget.Entry
}

type receiptRenderer struct {
	receiptsRepository port.ReceiptStorage
	accord             *widget.Accordion
	editPopup          *widget.PopUp
}

func (r *receiptRenderer) renderReceipts() *widget.Accordion {
	receipts := r.receiptsRepository.GetReceipts()
	aItems := make([]*widget.AccordionItem, 0, len(receipts))
	for _, receipt := range receipts {
		aItems = append(aItems, r.renderReceipt(receipt))
	}
	r.accord = widget.NewAccordion(aItems...)
	return r.accord
}

func (r *receiptRenderer) renderReceipt(receipt domain.Receipt) *widget.AccordionItem {
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
		itemsContainer.Add(container.New(
			layout.NewGridLayoutWithColumns(2),
			widget.NewLabel(item.Name),
			widget.NewLabel(item.Quantity),
		),
		)
	}
	receiptContainer.Add(itemsContainer)
	return widget.NewAccordionItem(receipt.Title, receiptContainer)
}

func (r *receiptRenderer) editForm() {

}
