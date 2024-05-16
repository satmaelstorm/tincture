package events

import "github.com/satmaelstorm/tincture/app/domain"

type TinctureDrunk struct {
	BaseEvent
	Tincture domain.Tincture
}

func (r *TinctureDrunk) Name() string {
	return "ready_tincture_drunk"
}

type TinctureBottled struct {
	BaseEvent
	Tincture domain.Tincture
}

func (r *TinctureBottled) Name() string {
	return "prepare_tincture_bottled"
}

type TinctureAddButton struct {
	BaseEvent
}

func (r *TinctureAddButton) Name() string {
	return "tincture_add_button"
}

type TinctureCancelButton struct {
	BaseEvent
}

func (r *TinctureCancelButton) Name() string {
	return "tincture_cancel_button"
}

type TinctureSubmit struct {
	BaseEvent
	Tincture domain.Tincture
}

func (r *TinctureSubmit) Name() string {
	return "tincture_submit"
}
