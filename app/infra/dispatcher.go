package infra

import "github.com/satmaelstorm/tincture/app/app/port"

type NoThreadSafeDispatcher struct {
	subscribers map[string][]port.Subscriber
}

func NewDispatcher() *NoThreadSafeDispatcher {
	d := new(NoThreadSafeDispatcher)
	d.subscribers = make(map[string][]port.Subscriber)
	return d
}

func (d *NoThreadSafeDispatcher) AddSubscriber(subscriber port.Subscriber) {
	for _, event := range subscriber.SupportEvents() {
		subscribers := d.subscribers[event.Name()]
		subscribers = append(subscribers, subscriber)
		d.subscribers[event.Name()] = subscribers
	}
}

func (d *NoThreadSafeDispatcher) Dispatch(event port.Event) {
	subscribers := d.subscribers[event.Name()]
	for _, subscriber := range subscribers {
		if event.IsPropagationStopped() {
			break
		}
		subscriber.DispatchEvent(event)
	}
}
