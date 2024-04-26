package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/satmaelstorm/tincture/app/app/port"
	"github.com/satmaelstorm/tincture/app/domain"
)

func renderReceipts(storage port.ReceiptStorage) []*widget.AccordionItem {
	receipts := storage.GetReceipts()
	aItems := make([]*widget.AccordionItem, 0, len(receipts))
	for _, receipt := range receipts {
		aItems = append(aItems, renderReceipt(receipt))
	}
	return aItems
}

func renderReceipt(receipt domain.Receipt) *widget.AccordionItem {
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
