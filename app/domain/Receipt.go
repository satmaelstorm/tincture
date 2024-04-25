package domain

import "github.com/google/uuid"

type Receipt struct {
	Uuid        uuid.UUID `gorm:"primarykey"`
	Title       string
	Description string
	Items       []ReceiptItem
}

type ReceiptItem struct {
	ID          uint `gorm:"primarykey"`
	ReceiptUuid uuid.UUID
	Name        string
	Quantity    string
}
