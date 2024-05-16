package handlers

import (
	"github.com/satmaelstorm/tincture/app/app/port"
	"github.com/satmaelstorm/tincture/app/app/renderers"
	"github.com/satmaelstorm/tincture/app/domain/events"
)

type ReadyTinctureHandlers struct {
	renderer   *renderers.ReadyTinctureRenderer
	repository port.TinctureStorage
}

func NewReadyTinctureHandlers(
	renderer *renderers.ReadyTinctureRenderer,
	repository port.TinctureStorage,
) *ReadyTinctureHandlers {
	return &ReadyTinctureHandlers{
		renderer:   renderer,
		repository: repository,
	}
}

func (r *ReadyTinctureHandlers) SupportEvents() []port.Event {
	return []port.Event{
		&events.TinctureDrunk{},
	}
}

func (r *ReadyTinctureHandlers) DispatchEvent(event port.Event) {
	switch e := event.(type) {
	case *events.TinctureDrunk:
		r.handleDeleteButton(e)
	}
}

func (r *ReadyTinctureHandlers) handleDeleteButton(event *events.TinctureDrunk) {
	r.repository.DeleteTincture(&event.Tincture)
	r.renderer.RemoveTincture(event.Tincture)
}
