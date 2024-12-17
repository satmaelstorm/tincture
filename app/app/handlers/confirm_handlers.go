package handlers

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/satmaelstorm/tincture/app/app/port"
	"github.com/satmaelstorm/tincture/app/domain/events"
)

type ConfirmFormHandler struct {
	w fyne.Window
}

func NewConfirmFormHandler(w fyne.Window) *ConfirmFormHandler {
	return &ConfirmFormHandler{w: w}
}

func (r *ConfirmFormHandler) SupportEvents() []port.Event {
	return []port.Event{
		&events.ReceiptConfirmDeleteButton{},
	}
}

func (r *ConfirmFormHandler) DispatchEvent(event port.Event) {
	switch e := event.(type) {
	case *events.ReceiptConfirmDeleteButton:
		r.confirmDeleteReceipt(e)
	}
}

func (r *ConfirmFormHandler) confirmDeleteReceipt(event *events.ReceiptConfirmDeleteButton) {
	dialog.ShowConfirm(
		"Вы действительно хотите удалить рецепт?",
		fmt.Sprintf("Вы хотите удалить рецепт %s", event.ReceiptTitle),
		event.Callback,
		r.w,
	)
}
