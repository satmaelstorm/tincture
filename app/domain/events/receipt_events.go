package events

import "github.com/google/uuid"

type ReceiptFormCancel struct {
	BaseEvent
}

func (r *ReceiptFormCancel) Name() string {
	return "receipt_form_cancel"
}

type ReceiptAddButton struct {
	BaseEvent
}

func (r *ReceiptAddButton) Name() string {
	return "receipt_add_button"
}

type ReceiptFormSubmit struct {
	BaseEvent
}

func (r *ReceiptFormSubmit) Name() string {
	return "receipt_from_submit"
}

type ReceiptEditButton struct {
	BaseEvent
	ReceiptUuid uuid.UUID
}

func (r *ReceiptEditButton) Name() string {
	return "receipt_edit_button"
}

type ReceiptDeleteButton struct {
	BaseEvent
	ReceiptUuid uuid.UUID
}

func (r *ReceiptDeleteButton) Name() string {
	return "receipt_delete_button"
}
