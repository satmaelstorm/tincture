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
	CreateReceipt(*domain.Receipt)
	SaveReceipt(*domain.Receipt)
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

type AppIcons interface {
}

type Event interface {
	Name() string
	IsPropagationStopped() bool
}

type Subscriber interface {
	DispatchEvent(Event)
	SupportEvents() []Event
}

type EventBus interface {
	Dispatch(Event)
}

type EventDispatcher interface {
	EventBus
	AddSubscriber(Subscriber)
}
