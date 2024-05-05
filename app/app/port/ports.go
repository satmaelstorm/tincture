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
	GetPreparingTinctures() []domain.Tincture
	GetReadyTinctures() []domain.Tincture
	SaveTincture(*domain.Tincture)
	CreateTincture(*domain.Tincture)
	DeleteTincture(*domain.Tincture)
}

type AppIcon interface {
	AsResource() fyne.Resource
}
