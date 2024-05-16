package handlers

import (
	"github.com/satmaelstorm/tincture/app/app/port"
	"github.com/satmaelstorm/tincture/app/app/renderers"
	"github.com/satmaelstorm/tincture/app/domain/events"
	"time"
)

type PrepareTinctureHandlers struct {
	ready   *renderers.ReadyTinctureRenderer
	prepare *renderers.PrepareTinctureRenderer

	repository port.TinctureStorage
}

func NewPrepareTinctureHandlers(
	ready *renderers.ReadyTinctureRenderer,
	prepare *renderers.PrepareTinctureRenderer,
	repository port.TinctureStorage,
) *PrepareTinctureHandlers {
	return &PrepareTinctureHandlers{
		ready:      ready,
		prepare:    prepare,
		repository: repository,
	}
}

func (p *PrepareTinctureHandlers) SupportEvents() []port.Event {
	return []port.Event{
		&events.TinctureBottled{},
		&events.TinctureSubmit{},
		&events.TinctureAddButton{},
		&events.TinctureCancelButton{},
	}
}

func (p *PrepareTinctureHandlers) DispatchEvent(event port.Event) {
	switch e := event.(type) {
	case *events.TinctureAddButton:
		p.handleAddButton()
	case *events.TinctureCancelButton:
		p.handleCancelButton()
	case *events.TinctureSubmit:
		p.handleSubmit(e)
	case *events.TinctureBottled:
		p.handleBottled(e)
	}
}

func (p *PrepareTinctureHandlers) handleBottled(event *events.TinctureBottled) {
	event.Tincture.Bottled(time.Now())
	p.prepare.RemoveTincture(event.Tincture)
	p.ready.AddTincture(event.Tincture)
}

func (p *PrepareTinctureHandlers) handleAddButton() {
	p.prepare.ShowAddPopup()
}

func (p *PrepareTinctureHandlers) handleCancelButton() {
	p.prepare.HideAddPopup()
}

func (p *PrepareTinctureHandlers) handleSubmit(event *events.TinctureSubmit) {
	p.repository.CreateTincture(&event.Tincture)
	p.prepare.AddTincture(event.Tincture)
	p.prepare.HideAddPopup()
}
