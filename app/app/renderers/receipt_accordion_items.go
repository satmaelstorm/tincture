package renderers

import (
	"fyne.io/fyne/v2/widget"
	"github.com/satmaelstorm/tincture/app/domain"
)

type receiptAccordionItems struct {
	items map[string]*widget.AccordionItem
}

func newReceiptAccordionItems(size int) *receiptAccordionItems {
	r := new(receiptAccordionItems)
	r.items = make(map[string]*widget.AccordionItem, size)
	return r
}

func (r *receiptAccordionItems) makeNew(receipt domain.Receipt) *widget.AccordionItem {
	r.items[receipt.Uuid.String()] = widget.NewAccordionItem("", widget.NewLabel(""))
	return r.items[receipt.Uuid.String()]
}

func (r *receiptAccordionItems) get(receipt domain.Receipt) *widget.AccordionItem {
	return r.items[receipt.Uuid.String()]
}
