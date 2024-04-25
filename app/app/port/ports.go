package port

import (
	"fyne.io/fyne/v2"
	"github.com/satmaelstorm/tincture/app/domain"
)

type InnerStorage interface {
	InitDb(fyne.URI) error
}

type ReceiptStorage interface {
	GetReceipts() []domain.Receipt
}

type TinctureStorage interface {
	GetTinctures() []domain.Tincture
}
