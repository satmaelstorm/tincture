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
	bus        port.EventBus
}

func NewReceiptFormHandlers(
	r *renderers.ReceiptRenderer,
	f *renderers.ReceiptEditForm,
	repo port.ReceiptStorage,
	bus port.EventBus,
) *ReceiptFormHandlers {
	return &ReceiptFormHandlers{renderer: r, form: f, repository: repo, bus: bus}
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
	switch e := event.(type) {
	case *events.ReceiptFormCancel:
		r.handleCancel()
	case *events.ReceiptAddButton:
		r.handleAddButton()
	case *events.ReceiptFormSubmit:
		r.handleSubmitButton()
	case *events.ReceiptDeleteButton:
		r.handleDeleteButton(e)
	case *events.ReceiptEditButton:
		r.handleEditButton(e)
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
		r.renderer.AddReceipt(receipt)
	} else {
		r.repository.SaveReceipt(&receipt)
		r.renderer.RefreshReceipt(receipt)
	}
	r.form.Hide()
}

func (r *ReceiptFormHandlers) handleDeleteButton(event *events.ReceiptDeleteButton) {
	receipt, ok := r.repository.GetReceipt(event.ReceiptUuid)
	if !ok {
		return
	}

	r.bus.Dispatch(&events.ReceiptConfirmDeleteButton{
		ReceiptTitle: receipt.Title,
		Callback: func(b bool) {
			if b && r.repository.DeleteReceipt(receipt) {
				r.renderer.RemoveReceipt(receipt)
			}
		},
	})
}

func (r *ReceiptFormHandlers) handleEditButton(event *events.ReceiptEditButton) {
	receipt, ok := r.repository.GetReceipt(event.ReceiptUuid)
	if !ok {
		r.handleAddButton()
		return
	}
	r.form.FromReceipt(receipt)
	r.form.Show()
}
