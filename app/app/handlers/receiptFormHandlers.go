package handlers

import (
	"github.com/satmaelstorm/tincture/app/app/port"
	"github.com/satmaelstorm/tincture/app/app/renderers"
	"github.com/satmaelstorm/tincture/app/domain/events"
)

type ReceiptFormHandlers struct {
	renderer   *renderers.ReceiptRenderer
	form       *renderers.ReceiptEditForm
	repository port.ReceiptStorage
}

func NewReceiptFormHandlers(
	r *renderers.ReceiptRenderer,
	f *renderers.ReceiptEditForm,
	repo port.ReceiptStorage,
) *ReceiptFormHandlers {
	return &ReceiptFormHandlers{renderer: r, form: f, repository: repo}
}

func (r *ReceiptFormHandlers) SupportEvents() []port.Event {
	return []port.Event{
		&events.ReceiptEditButton{},
		&events.ReceiptDeleteButton{},
		&events.ReceiptAddButton{},
		&events.ReceiptFormSubmit{},
		&events.ReceiptFormCancel{},
	}
}

func (r *ReceiptFormHandlers) DispatchEvent(event port.Event) {
	switch event.(type) {
	case *events.ReceiptFormCancel:
		r.handleCancel()
	case *events.ReceiptAddButton:
		r.handleAddButton()
	case *events.ReceiptFormSubmit:
		r.handleSubmitButton()
	case *events.ReceiptDeleteButton:

	}
}

func (r *ReceiptFormHandlers) handleCancel() {
	r.form.Clear()
	r.form.Hide()
}

func (r *ReceiptFormHandlers) handleAddButton() {
	r.form.Clear()
	r.form.Show()
}

func (r *ReceiptFormHandlers) handleSubmitButton() {
	receipt, isNew := r.form.CollectReceipt()
	if isNew {
		r.repository.CreateReceipt(&receipt)
	} else {
		r.repository.SaveReceipt(&receipt)
	}
	r.renderer.AddReceipt(receipt)
	r.form.Hide()
}

func (r *ReceiptFormHandlers) handleDeleteButton() {
	
}
