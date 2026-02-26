package handlers

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/satmaelstorm/tincture/app/app/port"
	"github.com/satmaelstorm/tincture/app/domain/events"
)

type ConfirmFormHandler struct {
	w         fyne.Window
	localizer port.Localizer
}

func NewConfirmFormHandler(w fyne.Window, localizer port.Localizer) *ConfirmFormHandler {
	return &ConfirmFormHandler{w: w, localizer: localizer}
}

func (r *ConfirmFormHandler) SupportEvents() []port.Event {
	return []port.Event{
		&events.ReceiptConfirmDeleteButton{},
		&events.TinctureConfirmDrunk{},
	}
}

func (r *ConfirmFormHandler) DispatchEvent(event port.Event) {
	switch e := event.(type) {
	case *events.ReceiptConfirmDeleteButton:
		r.confirmDeleteReceipt(e)
	case *events.TinctureConfirmDrunk:
		r.confirmDrunkTincture(e)

	}
}

func (r *ConfirmFormHandler) confirmDrunkTincture(event *events.TinctureConfirmDrunk) {
	dialog.ShowConfirm(
		r.localizer.Localize("tinctureConfirmDrink"),
		r.localizer.LocalizeWithTemplate("tinctureConfirmDelete", map[string]interface{}{
			"Name": event.Tincture.Name,
		}),
		event.Callback,
		r.w,
	)
}

func (r *ConfirmFormHandler) confirmDeleteReceipt(event *events.ReceiptConfirmDeleteButton) {
	dialog.ShowConfirm(
		r.localizer.Localize("receiptConfirmDelete"),
		r.localizer.LocalizeWithTemplate("receiptConfirmDeleteMessage", map[string]interface{}{
			"Title": event.ReceiptTitle,
		}),
		event.Callback,
		r.w,
	)
}
