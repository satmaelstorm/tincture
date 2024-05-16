package events

import "github.com/satmaelstorm/tincture/app/domain"

type TinctureDrunk struct {
	BaseEvent
	Tincture domain.Tincture
}

func (r *TinctureDrunk) Name() string {
	return "receipt_form_cancel"
}
