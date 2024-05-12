package domain

import "github.com/google/uuid"

type Receipt struct {
	Uuid        uuid.UUID `gorm:"primarykey"`
	Title       string
	Description string
	Items       []ReceiptItem
}

func (r *Receipt) Modify(
	title string,
	desc string,
	items ...ReceiptItem,
) {
	r.Title = title
	r.Description = desc
	r.Items = items
}

type ReceiptItem struct {
	ID          uint `gorm:"primarykey"`
	ReceiptUuid uuid.UUID
	Name        string
	Quantity    string
}

func NewReceipt(
	title, desc string,
	items ...ReceiptItem,
) Receipt {
	uuid := uuid.Must(uuid.NewV7())
	for i, _ := range items {
		items[i].ReceiptUuid = uuid
	}
	return Receipt{
		Uuid:        uuid,
		Title:       title,
		Description: desc,
		Items:       items,
	}
}

func NewReceiptItem(name, quantity string) ReceiptItem {
	return ReceiptItem{
		Name:     name,
		Quantity: quantity,
	}
}
